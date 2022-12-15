// SPDX-License-Identifier: GPL-3.0
pragma experimental ABIEncoderV2;

pragma solidity >=0.8.0 <0.9.0;

contract Aebmes {
    struct OWNER_DATA {
        address data_owner;
        mapping(bytes32 => string) wordsMap;
    }
    bytes16 private constant _SYMBOLS = "0123456789abcdef";

    struct TRAPDOOR_REQUEST {
        address data_user;
        address data_owner;
        bytes32 trapdoor_hash;
        uint256 endTimestamp;
    }

    struct FILE_METADATA {
        bytes32 fileHash;
    }

    mapping(address => OWNER_DATA) private ownerData;
    mapping(bytes32 => TRAPDOOR_REQUEST) private trapdoorList;
    mapping(string => bytes32) private fileDataList;

    // mapping(string => bytes32[]) private wordsMap;

    function addData(string[] memory words, string memory fileHash) public {
        OWNER_DATA storage owner_data = ownerData[msg.sender];
        if (owner_data.data_owner == address(0)) {
            owner_data.data_owner = msg.sender;
        }
        // bytes32 fileCipherTextHash = bytes32(fileCipherText);

        for (uint256 i = 0; i < words.length; i++) {
            bytes32 wordHash = createHash(words[i]);
            owner_data.wordsMap[wordHash] = fileHash;
        }
    }

    function viewData(address owner, string memory word)
        public
        view
        returns (string memory)
    {
        OWNER_DATA storage data = ownerData[owner];
        return data.wordsMap[createHash(word)];
    }

    function search_request(
        string[] memory words,
        address dataOwner,
        uint256 endTimestamp
    )
        public
        returns (
            address,
            string[] memory,
            bytes32
        )
    {
        string[] memory fileCipherTextData = new string[](words.length);
        //string[] memory b = new string[](size);
        fileCipherTextData[0] = "!";
        OWNER_DATA storage data = ownerData[dataOwner];
        uint64 latestIndex = 0;
        for (uint256 i = 0; i < words.length; i++) {
            bytes32 wordHash = createHash(words[i]);
            string memory cipherText = data.wordsMap[wordHash];
            bytes memory temp = bytes(cipherText);

            if (temp.length != 0) {
                fileCipherTextData[latestIndex] = cipherText;
                latestIndex++;
            }
        }
        bytes32 trapdoorHash = createTrapdoorHash(
            block.timestamp,
            endTimestamp,
            msg.sender
        );

        trapdoorList[trapdoorHash] = TRAPDOOR_REQUEST(
            msg.sender,
            dataOwner,
            trapdoorHash,
            endTimestamp
        );
        return (data.data_owner, fileCipherTextData, trapdoorHash);
    }

    // function viewTrapdoorRequest()
    function createHash(string memory data)
        private
        pure
        returns (bytes32 hashed_data)
    {
        hashed_data = keccak256(abi.encodePacked(data));
    }

    function createTrapdoorHash(
        uint256 t1,
        uint256 t2,
        address user_key
    ) internal pure returns (bytes32 trapdoorHash) {
        string[] memory stringArray = new string[](3);
        stringArray[0] = (toString(t1));
        stringArray[1] = (toString(t2));
        stringArray[2] = toString(uint256(uint160(user_key)));

        string memory finalString = concatenate(stringArray);
        trapdoorHash = createHash(finalString);
    }

    function toString(uint256 value) internal pure returns (string memory) {
        unchecked {
            uint256 length = log10(value) + 1;
            string memory buffer = new string(length);
            uint256 ptr;
            /// @solidity memory-safe-assembly
            assembly {
                ptr := add(buffer, add(32, length))
            }
            while (true) {
                ptr--;
                /// @solidity memory-safe-assembly
                assembly {
                    mstore8(ptr, byte(mod(value, 10), _SYMBOLS))
                }
                value /= 10;
                if (value == 0) break;
            }
            return buffer;
        }
    }

    function log10(uint256 value) internal pure returns (uint256) {
        uint256 result = 0;
        unchecked {
            if (value >= 10**64) {
                value /= 10**64;
                result += 64;
            }
            if (value >= 10**32) {
                value /= 10**32;
                result += 32;
            }
            if (value >= 10**16) {
                value /= 10**16;
                result += 16;
            }
            if (value >= 10**8) {
                value /= 10**8;
                result += 8;
            }
            if (value >= 10**4) {
                value /= 10**4;
                result += 4;
            }
            if (value >= 10**2) {
                value /= 10**2;
                result += 2;
            }
            if (value >= 10**1) {
                result += 1;
            }
        }
        return result;
    }

    function concatenate(string[] memory stringArray)
        internal
        pure
        returns (string memory stringValue)
    {
        for (uint256 j = 0; j < stringArray.length; j++) {
            stringValue = string(
                abi.encodePacked(stringValue, "", stringArray[j])
            );
        }
    }
}
