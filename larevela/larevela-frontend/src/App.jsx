import React, { useEffect, useState } from 'react';
import { Connection, Message, Transaction, clusterApiUrl } from '@solana/web3.js';
import PricingCards from './components/PricingCards';
import ContractLab from './pages/ContractLab';
import './App.css';

const DEVNET_USDC_MINT = '4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU';
const truncateAddress = (address) => {
  if (!address) {
    return '';
  }
  if (address.length <= 12) {
    return address;
  }
  return `${address.slice(0, 4)}...${address.slice(-4)}`;
};

const getSolflareProvider = () => {
  if (window.solflare?.isSolflare) {
    return window.solflare;
  }

  if (window.solflare && typeof window.solflare.connect === 'function') {
    return window.solflare;
  }

  if (window.solana?.isSolflare) {
    return window.solana;
  }

  if (Array.isArray(window.solana?.providers)) {
    return window.solana.providers.find((provider) => provider.isSolflare) || null;
  }

  if (Array.isArray(window.solflare?.providers)) {
    return window.solflare.providers.find((provider) => provider.isSolflare) || null;
  }

  return null;
};

function App() {
  const [pathname, setPathname] = useState(window.location.pathname || '/');
  const [walletAddress, setWalletAddress] = useState('');
  const [walletError, setWalletError] = useState('');
  const [isConnecting, setIsConnecting] = useState(false);
  const [isTrialSubmitting, setIsTrialSubmitting] = useState(false);
  const [trialMessage, setTrialMessage] = useState('');
  const [paymentAsset, setPaymentAsset] = useState('SOL');

  useEffect(() => {
    const provider = getSolflareProvider();
    if (!provider) {
      return undefined;
    }

    const syncAddressFromProvider = () => {
      const address = provider.publicKey?.toBase58?.() || '';
      setWalletAddress(address);
    };

    syncAddressFromProvider();

    const handleAccountChanged = (publicKey) => {
      setWalletAddress(publicKey?.toBase58?.() || '');
    };
    const handleDisconnect = () => {
      setWalletAddress('');
    };

    provider.on('accountChanged', handleAccountChanged);
    provider.on('disconnect', handleDisconnect);

    return () => {
      provider.removeListener('accountChanged', handleAccountChanged);
      provider.removeListener('disconnect', handleDisconnect);
    };
  }, []);

  useEffect(() => {
    const handlePopState = () => {
      setPathname(window.location.pathname || '/');
    };
    window.addEventListener('popstate', handlePopState);
    return () => {
      window.removeEventListener('popstate', handlePopState);
    };
  }, []);

  const navigateTo = (nextPath, state = {}) => {
    window.history.pushState(state, '', nextPath);
    setPathname(nextPath);
  };

  const handleStartFreeTrial = (plan) => {
    const decodeBase64ToBytes = (value) => Uint8Array.from(atob(value), (char) => char.charCodeAt(0));

    const submitTrial = async () => {
      setWalletError('');
      setTrialMessage('');
      setIsTrialSubmitting(true);

      try {
        const provider = getSolflareProvider();
        if (!provider) {
          throw new Error('Solflare not detected. Please install and unlock Solflare Wallet.');
        }

        let payerAddress = walletAddress;
        if (!payerAddress) {
          const connectResponse = await provider.connect();
          payerAddress = connectResponse?.publicKey?.toBase58?.() || provider.publicKey?.toBase58?.() || '';
          if (!payerAddress) {
            throw new Error('No account returned by Solflare.');
          }
          setWalletAddress(payerAddress);
        }

        const intentResp = await fetch('http://localhost:8888/api/v1/trade/payments/intents', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            orderNo: '',
            chainType: 'solana',
            network: 'devnet',
            chainId: 103,
            assetSymbol: paymentAsset,
            assetAddress: paymentAsset === 'USDC' ? DEVNET_USDC_MINT : '',
            payerAccount: payerAddress,
            referenceId: plan?.type || 'free_trial',
          }),
        });
        if (!intentResp.ok) {
          throw new Error('Failed to create payment intent.');
        }
        const intent = await intentResp.json();
        if (!intent?.serializedMessage || !intent?.paymentNo) {
          throw new Error('Payment intent is missing serializedMessage/paymentNo.');
        }
        localStorage.setItem('lastPaymentNo', intent.paymentNo);

        const messageBytes = decodeBase64ToBytes(intent.serializedMessage);
        const txMessage = Message.from(messageBytes);
        const unsignedTx = Transaction.populate(txMessage);

        let signature = '';
        if (provider.signAndSendTransaction) {
          const signed = await provider.signAndSendTransaction(unsignedTx);
          signature = signed?.signature || '';
        } else if (provider.signTransaction) {
          const signedTx = await provider.signTransaction(unsignedTx);
          const connection = new Connection(clusterApiUrl('devnet'), 'confirmed');
          signature = await connection.sendRawTransaction(signedTx.serialize());
        } else {
          throw new Error('Solflare provider does not support transaction signing.');
        }
        if (!signature) {
          throw new Error('No transaction signature returned by Solflare.');
        }

        const submitResp = await fetch('http://localhost:8888/api/v1/trade/payments/submit-tx', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            paymentNo: intent.paymentNo,
            txId: signature,
            fromAccount: payerAddress,
            fromTokenAccount: '',
          }),
        });
        if (!submitResp.ok) {
          throw new Error('Failed to submit transaction hash.');
        }
        localStorage.setItem('lastTxId', signature);

        setTrialMessage(`交易已提交并广播成功（${paymentAsset}）。 / Transaction signed, submitted, and broadcast successfully (${paymentAsset}).`);
        await new Promise((resolve) => setTimeout(resolve, 1500));
        navigateTo('/contract-lab', {});
      } catch (error) {
        setWalletError(error?.message || 'Start Free Trial failed.');
      } finally {
        setIsTrialSubmitting(false);
      }
    };

    setWalletError('');
    if (isTrialSubmitting) {
      return;
    }
    submitTrial();
  };

  const connectWallet = async () => {
    const provider = getSolflareProvider();
    if (!provider) {
      setWalletError('Solflare not detected. Please install and unlock Solflare Wallet.');
      return;
    }

    try {
      setWalletError('');
      setTrialMessage('');
      setIsConnecting(true);
      const response = await provider.connect();
      const address = response?.publicKey?.toBase58?.() || provider.publicKey?.toBase58?.() || '';
      if (!address) {
        throw new Error('No account returned by Solflare.');
      }
      setWalletAddress(address);
    } catch (_error) {
      setWalletError('Wallet connection failed or was rejected.');
    } finally {
      setIsConnecting(false);
    }
  };

  const disconnectWallet = async () => {
    const provider = getSolflareProvider();
    setWalletError('');
    setTrialMessage('');
    setIsConnecting(true);

    try {
      if (provider) {
        await provider.disconnect();
      }
    } catch (_error) {
      // Ignore disconnect failure and still clear local state.
    } finally {
      setWalletAddress('');
      setIsConnecting(false);
    }
  };

  const isContractLab = pathname === '/contract-lab';

  return (
    <div className="app">
      <header className="wallet-bar">
        <div className="wallet-status">
          {walletAddress ? (
            <span title={walletAddress}>Connected: {truncateAddress(walletAddress)}</span>
          ) : (
            'Wallet not connected'
          )}
        </div>
        <div className="payment-asset-selector">
          <label htmlFor="payment-asset">Select payment token</label>
          <select
            id="payment-asset"
            value={paymentAsset}
            onChange={(event) => setPaymentAsset(event.target.value)}
            disabled={isTrialSubmitting}
          >
            <option value="SOL">SOL</option>
            <option value="USDC">USDC</option>
          </select>
        </div>
        <button
          className="wallet-button"
          onClick={walletAddress ? disconnectWallet : connectWallet}
          disabled={isConnecting}
        >
          {isConnecting ? 'Processing...' : walletAddress ? 'Disconnect Wallet' : 'Connect Wallet'}
        </button>
      </header>
      {walletError && <p className="wallet-error">{walletError}</p>}
      {trialMessage && <p className="trial-message">{trialMessage}</p>}
      {!isContractLab ? (
        <>
          <PricingCards onStartFreeTrial={handleStartFreeTrial} isTrialSubmitting={isTrialSubmitting} />
        </>
      ) : (
        <ContractLab onBack={() => navigateTo('/', {})} />
      )}
    </div>
  );
}

export default App;