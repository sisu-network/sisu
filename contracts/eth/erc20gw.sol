pragma solidity >=0.8.4;

// OpenZeppelin's library: Context, Ownable
/**
 * @dev Provides information about the current execution context, including the
 * sender of the transaction and its data. While these are generally available
 * via msg.sender and msg.data, they should not be accessed in such a direct
 * manner, since when dealing with meta-transactions the account sending and
 * paying for execution may not be the actual sender (as far as an application
 * is concerned).
 *
 * This contract is only required for intermediate, library-like contracts.
 */
abstract contract Context {
    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }
}

/**
 * @dev Contract module which provides a basic access control mechanism, where
 * there is an account (an owner) that can be granted exclusive access to
 * specific functions.
 *
 * By default, the owner account will be the one that deploys the contract. This
 * can later be changed with {transferOwnership}.
 *
 * This module is used through inheritance. It will make available the modifier
 * `onlyOwner`, which can be applied to your functions to restrict their use to
 * the owner.
 */

abstract contract Ownable is Context {
    address private _owner;

    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

    /**
     * @dev Initializes the contract setting the deployer as the initial owner.
     */
    constructor() {
        _transferOwnership(_msgSender());
    }

    /**
     * @dev Returns the address of the current owner.
     */
    function owner() public view virtual returns (address) {
        return _owner;
    }

    /**
     * @dev Throws if called by any account other than the owner.
     */
    modifier onlyOwner() {
        require(owner() == _msgSender(), "Ownable: caller is not the owner");
        _;
    }

    /**
     * @dev Leaves the contract without owner. It will not be possible to call
     * `onlyOwner` functions anymore. Can only be called by the current owner.
     *
     * NOTE: Renouncing ownership will leave the contract without an owner,
     * thereby removing any functionality that is only available to the owner.
     */
    function renounceOwnership() public virtual onlyOwner {
        _transferOwnership(address(0));
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Can only be called by the current owner.
     */
    function transferOwnership(address newOwner) public virtual onlyOwner {
        require(newOwner != address(0), "Ownable: new owner is the zero address");
        _transferOwnership(newOwner);
    }

    /**
     * @dev Transfers ownership of the contract to a new account (`newOwner`).
     * Internal function without access restriction.
     */
    function _transferOwnership(address newOwner) internal virtual {
        address oldOwner = _owner;
        _owner = newOwner;
        emit OwnershipTransferred(oldOwner, newOwner);
    }
}

// helper methods for interacting with ERC20 tokens and sending NATIVE that do not consistently return true/false
library TransferHelper {
    function safeApprove(address token, address to, uint value) internal {
        // bytes4(keccak256(bytes('approve(address,uint256)')));
        (bool success, bytes memory data) = token.call(abi.encodeWithSelector(0x095ea7b3, to, value));
        require(success && (data.length == 0 || abi.decode(data, (bool))), 'TransferHelper: APPROVE_FAILED');
    }

    function safeTransfer(address token, address to, uint value) internal {
        // bytes4(keccak256(bytes('transfer(address,uint256)')));
        (bool success, bytes memory data) = token.call(abi.encodeWithSelector(0xa9059cbb, to, value));
        require(success && (data.length == 0 || abi.decode(data, (bool))), 'TransferHelper: TRANSFER_FAILED');
    }

    function safeTransferFrom(address token, address from, address to, uint value) internal {
        // bytes4(keccak256(bytes('transferFrom(address,address,uint256)')));
        (bool success, bytes memory data) = token.call(abi.encodeWithSelector(0x23b872dd, from, to, value));
        require(success && (data.length == 0 || abi.decode(data, (bool))), 'TransferHelper: TRANSFER_FROM_FAILED');
    }

    function safeTransferNative(address to, uint value) internal {
        (bool success,) = to.call{value : value}(new bytes(0));
        require(success, 'TransferHelper: NATIVE_TRANSFER_FAILED');
    }
}

/**
 * @dev Interface of the ERC20 standard as defined in the EIP.
 */
interface IERC20 {
    function totalSupply() external view returns (uint256);

    function balanceOf(address account) external view returns (uint256);

    function transfer(address recipient, uint256 amount) external returns (bool);

    function allowance(address owner, address spender) external view returns (uint256);

    function approve(address spender, uint256 amount) external returns (bool);

    function permit(address target, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) external;

    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    function transferWithPermit(address target, address to, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

contract ERC20Gateway is Ownable {
    bool public pause;

    // list of supported chains
    // key: chainId
    // value: true if this chain is supported
    mapping(string => bool) public supportedChains;

    constructor(string[] memory _supportedChains) {
        pause = false;

        for (uint i = 0; i < _supportedChains.length; i++) {
            supportedChains[_supportedChains[i]] = true;
        }
    }

    event TransferOutEvent(string indexed destChain, address indexed token, address indexed sender, uint256 amount);
    event TransferInEvent(address indexed token, address indexed reipient, uint256 amount);
    event RemoveSupportedChainEvent(string indexed chain);
    event AddSupportedChainEvent(string indexed chain);

    modifier isNotPaused {
        require(pause == false, "Gateway is paused");
        _;
    }

    // User can call TransferOut to deposit their ERC20 token to gateway
    // Anyone can call TransferOut
    function TransferOut(string memory destChain, address _token, uint256 _amount) public isNotPaused {
        require(supportedChains[destChain] == true, "destChain is not supported");
        TransferHelper.safeTransferFrom(_token, msg.sender, address(this), _amount);

        emit TransferOutEvent(destChain, _token, msg.sender, _amount);
    }

    // Pool owner call TransferIn to release user's ERC20 token in destination chain
    // Triggered by bridge's backend
    function TransferIn(address _token, address recipient, uint256 _amount) public onlyOwner isNotPaused {
        uint256 gwBalance = IERC20(_token).balanceOf(address(this));
        require(gwBalance >= _amount, "Gateway balance is less than required amount");

        TransferHelper.safeTransferFrom(_token, address(this), recipient, _amount);

        emit TransferInEvent(_token, recipient, _amount);
    }

    function PauseGateway() public onlyOwner {
        require(pause == false, "Gateway is paused already");
        pause = true;
    }

    function ResumeGateway() public onlyOwner {
        require(pause == true, "Gateway is not paused already");
        pause = false;
    }

    function RemoveSupportedChain(string memory chain) public onlyOwner {
        supportedChains[chain] = false;

        emit RemoveSupportedChainEvent(chain);
    }

    function AddSupportedChain(string memory chain) public onlyOwner {
        supportedChains[chain] = true;

        emit AddSupportedChainEvent(chain);
    }
}
