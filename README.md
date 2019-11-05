# MPSE-Blockchain-P2P-PBFT

required:
    
*   openconnect for vpn and URL firewall.fbi.h-da.de

*  golibp2p library 

*  golang 

1. P2P- Nodes:
    * VPN to hda required (for Linux use openconnect and use url firewall.fbi.h-da.de) login with your ist... account
    * `sudo openconnect firewall.fbi.h-da.de`
    * required the public and private key stored in vm_sshkey (key from Anton)
    * `sudo cp vm_sshkey/id_rsa ~/.ssh/`
    * `chmod 400 ~/.ssh/id_rsa` 
    * using the URLs:root@uvm-isa-tokenchain-1 (1 bis 40) with ssh
    * `ssh root@uvm-isa-tokenchain-(1..40)`
    * clone the p2p file in every node and execute it

This project is just to create a p2p network and pbft algorithm for synchronization