# ubuntuCveExtractor

This is an example project to try out crawler in golang.

The goal is to extract CVE information for Ubuntu releases. 

It will read data start in the page `https://ubuntu.com/security/cve`
and will extract usefull information for each CVE present in this list.

Data is saved in a json file name `ubuntu_cve.json`

## Output Data Format
