// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface IERC20 {
    function totalSupply() external view returns (uint256);

    function decimals() external view returns (uint8);

    function balanceOf(address account) external view returns (uint256);

    function transfer(address recipient, uint256 amount)
        external
        returns (bool);

    function allowance(address owner, address spender)
        external
        view
        returns (uint256);

    function approve(address spender, uint256 amount) external returns (bool);

    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(
        address indexed owner,
        address indexed spender,
        uint256 value
    );
}

library Address {
    // https://soliditydeveloper.com/extcodehash
    function isContract(address account) internal view returns (bool) {
        bytes32 codehash;
        bytes32 accountHash = 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470;
        // solhint-disable-next-line no-inline-assembly
        assembly {
            codehash := extcodehash(account)
        }
        return (codehash != 0x0 && codehash != accountHash);
    }
}

library TransferHelper {
    function safeTransferNative(address to, uint256 value) internal {
        (bool success, ) = to.call{value: value}(new bytes(0));
        require(success, "TransferHelper: NATIVE_TRANSFER_FAILED");
    }
}

/**
 * @title SafeERC20
 * @dev Wrappers around ERC20 operations that throw on failure (when the token
 * contract returns false). Tokens that return no value (and instead revert or
 * throw on failure) are also supported, non-reverting calls are assumed to be
 * successful.
 * To use this library you can add a `using SafeERC20 for IERC20;` statement to your contract,
 * which allows you to call the safe operations as `token.safeTransfer(...)`, etc.
 */
library SafeERC20 {
    using Address for address;

    function safeTransfer(
        IERC20 token,
        address to,
        uint256 value
    ) internal {
        callOptionalReturn(
            token,
            abi.encodeWithSelector(token.transfer.selector, to, value)
        );
    }

    function safeTransferFrom(
        IERC20 token,
        address from,
        address to,
        uint256 value
    ) internal {
        callOptionalReturn(
            token,
            abi.encodeWithSelector(token.transferFrom.selector, from, to, value)
        );
    }

    /**
     * @dev Imitates a Solidity high-level call (i.e. a regular function call to a contract), relaxing the requirement
     * on the return value: the return value is optional (but if data is returned, it must not be false).
     * @param token The token targeted by the call.
     * @param data The call data (encoded using abi.encode or one of its variants).
     */
    function callOptionalReturn(IERC20 token, bytes memory data) private {
        require(address(token).isContract(), "SafeERC20: call to non-contract");

        // solhint-disable-next-line avoid-low-level-calls
        (bool success, bytes memory returndata) = address(token).call(data);
        require(success, "SafeERC20: low-level call failed");

        if (returndata.length > 0) {
            // Return data is optional
            // solhint-disable-next-line max-line-length
            require(
                abi.decode(returndata, (bool)),
                "SafeERC20: ERC20 operation did not succeed"
            );
        }
    }
}

contract Vault {
    using SafeERC20 for IERC20;

    // Symbolic spender address in the balance map. Spenders could change but this address remains
    // the same.
    address private constant symbolicSpender =
        0x3000000000000000000000000000000000000000;
    // Symbolic address of native token
    address private constant native =
        0x4000000000000000000000000000000000000000;

    mapping(address => bool) spenders;
    address private admin;

    mapping(address => mapping(address => uint256)) private balances;

    modifier onlySpender() {
        require(spenders[msg.sender], "Not spender: FORBIDDEN");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Not admin: FORBIDDEN");
        _;
    }

    // The vault does not have enough balance to transfer token to recipient. Temporarily increases
    // user's balance for later withdrawal.
    event Code501();

    // Retry transfer fails.
    event Code502();

    constructor() {
        admin = msg.sender;
    }

    ////////////////////////////////////////////////////////////////
    // Admin
    ////////////////////////////////////////////////////////////////

    function addSpender(address spender) external onlyAdmin {
        spenders[spender] = true;
    }

    function removeSpender(address spender) external onlyAdmin {
        spenders[spender] = false;
    }

    function changeAdmin(address newAdmin) external onlyAdmin {
        admin = newAdmin;
    }

    ////////////////////////////////////////////////////////////////
    // Deposit/Withdraws
    ////////////////////////////////////////////////////////////////

    /**
     * @dev Deposits a certain amount of token into this liquidity.
     */
    function deposit(address token, uint256 amount) external {
        _deposit(token, msg.sender, amount);
    }

    /**
     * @dev Deposits a certain amount of token into symbolic spender account.
     */
    function depositFor(
        address token,
        address receiver,
        uint256 amount
    ) external {
        IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
        _inc(token, receiver, amount);
    }

    /**
     * @dev Deposits native token into this vault.
     */
    function depositNative() external payable {
        _inc(native, msg.sender, msg.value);
    }

    /**
     * @dev Deposits native token into this vault for a receiver.
     */
    function depositNativeFor(address receiver) external payable {
        _inc(native, receiver, msg.value);
    }

    /**
     * @dev Withdraws a certain amount of token from sender's account to a `to` account. This
     * function transfers assets from this liquidity to the receiver's address.
     */
    function withdraw(
        address token,
        address to,
        uint256 amount
    ) external {
        _withdraw(token, msg.sender, to, amount);
    }

    /**
     * @dev Withdraw native token to a `to` address.
     */
    function withdrawNative(address to, uint256 amount) external {
        _dec(native, msg.sender, amount);
        TransferHelper.safeTransferNative(to, amount);
    }

    /**
     * @dev Withdraw native token to a `to` address.
     */
    function transferInNative(address to, uint256 amount) external onlySpender {
        if (address(this).balance >= amount) {
            TransferHelper.safeTransferNative(to, amount);
        } else {
            _inc(native, to, amount);
            emit Code501();
        }
    }

    function retryTransferNative(address to, uint256 amount)
        external
        onlySpender
    {
        uint256 actual = amount;
        if (actual > balances[native][to]) {
            actual = balances[native][to];
        }

        if (address(this).balance >= actual) {
            TransferHelper.safeTransferNative(to, actual);
        } else {
            emit Code502();
        }
    }

    /**
     * @dev Transfer an `amount` of ERC20 token to a `to` address from this contract's account.
     */
    function transferOutNonEvm(
        address token,
        string memory dstChain,
        string memory to,
        uint256 amount
    ) public {
        if (balances[token][msg.sender] >= amount) {
            _dec(token, msg.sender, amount);
        } else {
            IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
        }
    }

    /**
     * @dev Transfer an `amount` of ERC20 token to a `to` address from this contract's account. The
     *   `to` recipient is an address instead of a string.
     */
    function transferOut(
        address token,
        string memory dstChain,
        address to,
        uint256 amount
    ) public {
        if (balances[token][msg.sender] >= amount) {
            _dec(token, msg.sender, amount);
        } else {
            IERC20(token).safeTransferFrom(msg.sender, address(this), amount);
        }
    }

    /**
     * @dev Transfer multiple tokens out to different destination.
     */
    function transferOutMultipleNonEvm(
        address[] memory tokens,
        string[] memory dstChains,
        string[] memory tos,
        uint256[] memory amounts
    ) external {
        for (uint32 i = 0; i < tokens.length; i++) {
            transferOutNonEvm(tokens[i], dstChains[i], tos[i], amounts[i]);
        }
    }

    /**
     * @dev Transfer multiple tokens out to different destination.The `to` recipient is an address
     * instead of a string.
     */
    function transferOutMultiple(
        address[] memory tokens,
        string[] memory dstChains,
        address[] memory tos,
        uint256[] memory amounts
    ) external {
        for (uint32 i = 0; i < tokens.length; i++) {
            transferOut(tokens[i], dstChains[i], tos[i], amounts[i]);
        }
    }

    /**
     * @dev Transfer out native token from an account to a new chain.
     */
    function transferOutNative(string memory to, string memory dstChain)
        external
        payable
    {
        _inc(native, msg.sender, msg.value);
    }

    /**
     * @dev Transfer an `amount` of token to a `to` address from this contract's account.
     */
    function transferIn(
        address token,
        address to,
        uint256 amount
    ) public onlySpender {
        if (IERC20(token).balanceOf(address(this)) >= amount) {
            IERC20(token).safeTransfer(to, amount);
        } else {
            _inc(token, to, amount);
            emit Code501();
        }
    }

    /**
     * @dev Transfer an `amount` of token to a `to` address from this contract's account.
     */
    function transferInMultiple(
        address[] memory tokens,
        address[] memory tos,
        uint256[] memory amounts
    ) external onlySpender {
        for (uint32 i = 0; i < tokens.length; i++) {
            transferIn(tokens[i], tos[i], amounts[i]);
        }
    }

    /**
     * @dev Retries failed transfer transaction. It's possible that the amount has been
     */
    function retryTransfer(
        address token,
        address to,
        uint256 amount
    ) external onlySpender {
        uint256 actual = amount;
        if (actual > balances[token][to]) {
            actual = balances[token][to];
        }

        if (IERC20(token).balanceOf(address(this)) >= actual) {
            _withdraw(token, to, to, actual);
        } else {
            emit Code502();
        }
    }

    function _deposit(
        address token,
        address from,
        uint256 amount
    ) internal {
        IERC20(token).safeTransferFrom(from, address(this), amount);

        _inc(token, from, amount);
    }

    function _withdraw(
        address token,
        address from,
        address to,
        uint256 amount
    ) internal {
        _dec(token, from, amount);
        IERC20(token).safeTransfer(to, amount);
    }

    function _dec(
        address token,
        address account,
        uint256 amount
    ) internal {
        require(account != address(0), "dec: address is 0");
        require(
            balances[token][account] >= amount,
            "dec: amount exceeds balance"
        );

        balances[token][account] -= amount;
    }

    function _inc(
        address token,
        address acc,
        uint256 amount
    ) internal {
        require(acc != address(0), "inc: address is 0");

        balances[token][acc] += amount;
    }

    /**
     * @dev Returns token balance of an account.
     */
    function balanceOf(address token, address account)
        external
        view
        returns (uint256)
    {
        return balances[token][account];
    }
}
