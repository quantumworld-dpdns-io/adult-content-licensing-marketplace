// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract LicenseRegistry {
    struct License {
        string creatorId;
        string title;
        bool aiTrainingProhibited;
    }

    mapping(string => License) private licenses;

    event LicenseCreated(string indexed licenseId, string creatorId, string title);

    function createLicense(string calldata licenseId, string calldata creatorId, string calldata title, bool aiTrainingProhibited) external {
        require(bytes(licenses[licenseId].title).length == 0, "exists");
        licenses[licenseId] = License({
            creatorId: creatorId,
            title: title,
            aiTrainingProhibited: aiTrainingProhibited
        });
        emit LicenseCreated(licenseId, creatorId, title);
    }
}
