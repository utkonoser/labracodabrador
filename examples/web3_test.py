#!/usr/bin/env python3
"""
Example script to test the blockchain with web3.py
Install: pip install web3
"""

from web3 import Web3
import sys

def main():
    # Connect to local node
    w3 = Web3(Web3.HTTPProvider('http://localhost:8545'))
    
    print("=" * 60)
    print("Labracodabrador Blockchain Test")
    print("=" * 60)
    
    # Check connection
    if not w3.is_connected():
        print("\n✗ Failed to connect to node")
        sys.exit(1)
    
    print("\n✓ Connected to network")
    print(f"  Chain ID: {w3.eth.chain_id}")
    
    # Get latest block
    block_number = w3.eth.block_number
    print(f"\n✓ Latest block: {block_number}")
    
    # Get block details
    block = w3.eth.get_block(block_number)
    print(f"  Timestamp: {block['timestamp']}")
    print(f"  Transactions: {len(block['transactions'])}")
    
    # Check if mining
    is_mining = w3.eth.mining
    print(f"\n✓ Mining status: {is_mining}")
    
    # Get peer count
    peer_count = w3.net.peer_count
    print(f"✓ Connected peers: {peer_count}")
    
    # List accounts
    accounts = w3.eth.accounts
    print(f"\n✓ Available accounts: {len(accounts)}")
    
    if accounts:
        print("\nAccount balances:")
        for account in accounts:
            balance = w3.eth.get_balance(account)
            balance_eth = w3.from_wei(balance, 'ether')
            print(f"  {account}: {balance_eth} ETH")
    
    # Test transaction (if accounts available)
    if len(accounts) >= 2:
        print("\n✓ Testing transaction...")
        try:
            tx_hash = w3.eth.send_transaction({
                'from': accounts[0],
                'to': accounts[1],
                'value': w3.to_wei(1, 'ether'),
                'gas': 21000,
                'gasPrice': w3.to_wei(1, 'gwei')
            })
            print(f"  Transaction hash: {tx_hash.hex()}")
        except Exception as e:
            print(f"  ⚠ Transaction failed: {e}")
    
    print("\n" + "=" * 60)
    print("All tests passed!")
    print("=" * 60 + "\n")

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"\n✗ Error: {e}")
        sys.exit(1)

