// Example script to test the blockchain with ethers.js
// Install: npm install ethers

const { ethers } = require("ethers");

async function main() {
    // Connect to local node
    const provider = new ethers.JsonRpcProvider("http://localhost:8545");
    
    console.log("=".repeat(60));
    console.log("Labracodabrador Blockchain Test");
    console.log("=".repeat(60));
    
    // Check connection
    try {
        const network = await provider.getNetwork();
        console.log("\n✓ Connected to network");
        console.log("  Chain ID:", network.chainId.toString());
        
        // Get latest block
        const blockNumber = await provider.getBlockNumber();
        console.log("\n✓ Latest block:", blockNumber);
        
        // Get block details
        const block = await provider.getBlock(blockNumber);
        console.log("  Timestamp:", new Date(block.timestamp * 1000).toISOString());
        console.log("  Transactions:", block.transactions.length);
        
        // Check if mining
        const isMining = await provider.send("eth_mining", []);
        console.log("\n✓ Mining status:", isMining);
        
        // Get peer count
        const peerCount = await provider.send("net_peerCount", []);
        console.log("✓ Connected peers:", parseInt(peerCount, 16));
        
        // List accounts
        const accounts = await provider.send("eth_accounts", []);
        console.log("\n✓ Available accounts:", accounts.length);
        
        if (accounts.length > 0) {
            console.log("\nAccount balances:");
            for (const account of accounts) {
                const balance = await provider.getBalance(account);
                console.log(`  ${account}: ${ethers.formatEther(balance)} ETH`);
            }
        }
        
        console.log("\n" + "=".repeat(60));
        console.log("All tests passed!");
        console.log("=".repeat(60) + "\n");
        
    } catch (error) {
        console.error("\n✗ Error:", error.message);
        process.exit(1);
    }
}

main()
    .then(() => process.exit(0))
    .catch(error => {
        console.error(error);
        process.exit(1);
    });

