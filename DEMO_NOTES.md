# Demo Notes

Notes, tips, and hints for using the various Veracode scan types with this application.

Also see the `docs/flaws` folder for in-depth explanations of the various exploits exposed in this application.

## Static scanning

Build the app:

	<run command here>

The file is the file to upload for scanning. Either upload this file to the Veracode platform for a Policy/Sandbox scan, or use it with the Veracode Pipeline scan.

## SCA scanning

### Upload/Scan

This will just happen as part of the Policy or Sandbox scan.

### Agent-based scan

Use either the command-line version of the SCA agent (follow the install and config instructions in the Veracode [docs center](https://docs.veracode.com/r/c_sc_what_is) ) or the IDE plugin to initiate an SCA scan.

### Vulnerable Methods

The Veracode agent-based SCA scan can also find [vulnerable methods](https://docs.veracode.com/r/Finding_and_Fixing_Vulnerabilities#fixing-vulnerable-methods).  Currently no vulnerable methods are implemented in VerademoGO though.