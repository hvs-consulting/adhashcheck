# ADHashCheck

Compare the NTLM hashes of all Active Directory accounts to find accounts with same password

The tool can read the NTLM hashes either from a secretsdump, from a CSV file or from the NTDS.dit and SYSTEM registry hive.

# Warning

If you follow this guide, you extract all NTLM hashes from your Active Directory. These hashes can be used for authentication, so if someone gets access to it, **your Active Directory is fully compromised**. You are responsible to verify the integrity of the used third-party libraries and tools.

HvS-Consulting AG is not liably for any damage caused by the use of this tool.

Security tools also do not like this.

# Use it

## Get NTLM hashes

### Option 1: Get ntds.dit and system registriy hive from DC

Export the system registry hive:

    reg save HKLM\SYSTEM system.hive

Get the ntds.dit (for example with FTKImager)

### Option 2: DCSync

The hashes can be collected with DCSync using [impacket's](https://github.com/SecureAuthCorp/impacket) secretsdump:

    python secretsdump.py -just-dc-ntlm <domain>/<user>[:<password>]@<dc> > secretsdump.txt

### Option 3: CSV file

Got the hashes somehow different? The tool also accepts a CSV file in the format:

    username,hash

## Analyze

Depending on the input files, run:

For the ntds.dit option:

    ./adhashcheck --ntds <ntds.dit> --system <system.hive> --output <output directory>

For the secretsdump option:

    ./adhashcheck --secretsdump <secretsdump.txt> --output <output directory>

For the CSV option:

    ./adhashcheck --csv <hashes.csv> --output <output directory>


The output directory now contains three files:

* `reuse.csv`: Password hashes that are used more than once (including count and usernames)
* `reuse-without-hash.csv`: Like `reuse.csv`, but without NTLM hashes


# License

Copyright 2021 HvS-Consulting AG

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
