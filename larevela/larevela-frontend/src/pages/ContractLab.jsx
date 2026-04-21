import React, { useState } from 'react';
import './ContractLab.css';

const getMetaMaskProvider = () => {
  const { ethereum } = window;
  if (!ethereum) {
    return null;
  }
  if (Array.isArray(ethereum.providers) && ethereum.providers.length > 0) {
    return ethereum.providers.find((provider) => provider.isMetaMask) || null;
  }
  return ethereum.isMetaMask ? ethereum : null;
};

const formatEthFromHexWei = (hexWei) => {
  try {
    const wei = BigInt(hexWei);
    const whole = wei / 1000000000000000000n;
    const fraction = (wei % 1000000000000000000n).toString().padStart(18, '0').slice(0, 6);
    return `${whole.toString()}.${fraction} ETH`;
  } catch (_error) {
    return '0 ETH';
  }
};

const parseEthToHexWei = (ethValue) => {
  const value = ethValue.trim();
  if (!/^\d+(\.\d+)?$/.test(value)) {
    throw new Error('Invalid ETH amount format.');
  }
  const [wholePart, fractionPart = ''] = value.split('.');
  const safeFraction = fractionPart.slice(0, 18).padEnd(18, '0');
  const wei = BigInt(wholePart) * 1000000000000000000n + BigInt(safeFraction || '0');
  if (wei <= 0n) {
    throw new Error('Amount must be greater than 0.');
  }
  return `0x${wei.toString(16)}`;
};

const getExplorerTxUrl = (chainId, txHash) => {
  if (chainId === '0xaa36a7') {
    return `https://sepolia.etherscan.io/tx/${txHash}`;
  }
  if (chainId === '0x1') {
    return `https://etherscan.io/tx/${txHash}`;
  }
  return '';
};

const ContractLab = ({ onBack }) => {
  const [readResult, setReadResult] = useState('');
  const [readLoading, setReadLoading] = useState(false);
  const [recipientAddress, setRecipientAddress] = useState('');
  const [amount, setAmount] = useState('0.0001');
  const [txStatus, setTxStatus] = useState('idle');
  const [txMessage, setTxMessage] = useState('');
  const [txHash, setTxHash] = useState('');
  const [txExplorerUrl, setTxExplorerUrl] = useState('');

  const ensureWalletConnected = async () => {
    const provider = getMetaMaskProvider();
    if (!provider) {
      throw new Error('MetaMask not found. Please install and unlock MetaMask.');
    }

    const accounts = await provider.request({ method: 'eth_accounts' });
    if (accounts?.length > 0) {
      return { provider, address: accounts[0] };
    }

    const requestedAccounts = await provider.request({ method: 'eth_requestAccounts' });
    const address = requestedAccounts?.[0];
    if (!address) {
      throw new Error('No account returned by MetaMask.');
    }
    return { provider, address };
  };

  const readOnchainData = async () => {
    setReadLoading(true);
    setReadResult('');

    try {
      const { provider, address } = await ensureWalletConnected();
      const [balanceHex, chainId] = await Promise.all([
        provider.request({ method: 'eth_getBalance', params: [address, 'latest'] }),
        provider.request({ method: 'eth_chainId' }),
      ]);

      setReadResult(`Address: ${address} | Balance: ${formatEthFromHexWei(balanceHex)} | Chain ID: ${chainId}`);
    } catch (error) {
      setReadResult(error?.message || 'Read failed. Please retry.');
    } finally {
      setReadLoading(false);
    }
  };

  const submitWriteTx = async () => {
    try {
      const { provider, address } = await ensureWalletConnected();
      const value = parseEthToHexWei(amount);
      const to = recipientAddress || address;
      const chainId = await provider.request({ method: 'eth_chainId' });

      setTxStatus('pending');
      setTxMessage('Transaction submitted, waiting for confirmation...');
      setTxHash('');
      setTxExplorerUrl('');

      const hash = await provider.request({
        method: 'eth_sendTransaction',
        params: [
          {
            from: address,
            to,
            value,
          },
        ],
      });
      setTxHash(hash);
      setTxExplorerUrl(getExplorerTxUrl(chainId, hash));

      let receipt = null;
      for (let i = 0; i < 30; i += 1) {
        // Poll receipt until tx gets mined.
        // eslint-disable-next-line no-await-in-loop
        receipt = await provider.request({
          method: 'eth_getTransactionReceipt',
          params: [hash],
        });
        if (receipt) {
          break;
        }
        // eslint-disable-next-line no-await-in-loop
        await new Promise((resolve) => setTimeout(resolve, 1500));
      }

      if (!receipt) {
        setTxStatus('failed');
        setTxMessage('Timed out waiting for on-chain confirmation.');
        return;
      }

      if (receipt.status === '0x1') {
        setTxStatus('success');
        setTxMessage('Transaction confirmed on-chain.');
      } else {
        setTxStatus('failed');
        setTxMessage('Transaction failed on-chain.');
      }

      if (!getExplorerTxUrl(chainId, hash)) {
        setTxMessage((prev) => `${prev} (Explorer link unavailable for current network.)`);
      }
    } catch (error) {
      setTxStatus('failed');
      setTxMessage(error?.message || 'Transaction failed or was rejected.');
    }
  };

  return (
    <section className="contract-lab">
      <div className="lab-header">
        <button className="lab-back" onClick={onBack}>
          Back to Pricing
        </button>
        <h2>Contract Interaction Lab (MetaMask Only)</h2>
      </div>

      <div className="lab-card">
        <h3>1) Read (View/Pure style call)</h3>
        <p className="lab-muted">Read-only RPC calls: fetch account ETH balance and current chain id.</p>
        <button className="lab-btn" onClick={readOnchainData} disabled={readLoading}>
          {readLoading ? 'Reading...' : 'Read On-chain Data'}
        </button>
        {readResult && <p className="lab-result">{readResult}</p>}
      </div>

      <div className="lab-card">
        <h3>2) Write (State-changing tx)</h3>
        <p className="lab-muted">Send ETH transfer with MetaMask and track transaction status.</p>
        <label className="lab-label" htmlFor="recipient">
          Recipient address (optional, defaults to your connected wallet)
        </label>
        <input
          id="recipient"
          className="lab-input"
          placeholder="Enter EVM address (0x...)"
          value={recipientAddress}
          onChange={(event) => setRecipientAddress(event.target.value.trim())}
        />
        <label className="lab-label" htmlFor="amount">
          Amount (ETH)
        </label>
        <input
          id="amount"
          className="lab-input"
          value={amount}
          onChange={(event) => setAmount(event.target.value)}
        />
        <button className="lab-btn" onClick={submitWriteTx} disabled={txStatus === 'pending'}>
          {txStatus === 'pending' ? 'Pending...' : 'Send Transaction'}
        </button>
        <p className={`lab-status ${txStatus}`}>Status: {txStatus === 'idle' ? 'Not started' : txStatus}</p>
        {txMessage && <p className="lab-result">{txMessage}</p>}
        {txHash && txExplorerUrl && (
          <a
            className="lab-link"
            href={txExplorerUrl}
            target="_blank"
            rel="noreferrer"
          >
            View transaction on Explorer
          </a>
        )}
      </div>
    </section>
  );
};

export default ContractLab;
