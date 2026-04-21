import React, { useEffect, useMemo, useState } from 'react';
import PricingCards from './components/PricingCards';
import ContractLab from './pages/ContractLab';
import './App.css';

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

function App() {
  const [pathname, setPathname] = useState(window.location.pathname || '/');
  const [walletAddress, setWalletAddress] = useState('');
  const [walletError, setWalletError] = useState('');
  const [isConnecting, setIsConnecting] = useState(false);

  const shortAddress = useMemo(() => {
    if (!walletAddress) {
      return '';
    }
    return `${walletAddress.slice(0, 6)}...${walletAddress.slice(-4)}`;
  }, [walletAddress]);

  useEffect(() => {
    const provider = getMetaMaskProvider();
    if (!provider) {
      return undefined;
    }

    const initWallet = async () => {
      try {
        const accounts = await provider.request({ method: 'eth_accounts' });
        if (accounts?.length > 0) {
          setWalletAddress(accounts[0]);
        }
      } catch (_error) {
        // ignore silent init error
      }
    };

    initWallet();

    const handleAccountChanged = (accounts) => {
      setWalletAddress(accounts?.[0] || '');
    };
    provider.on('accountsChanged', handleAccountChanged);

    return () => {
      provider.removeListener('accountsChanged', handleAccountChanged);
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

  const handleStartFreeTrial = () => {
    navigateTo('/contract-lab', {});
  };

  const connectWallet = async () => {
    const provider = getMetaMaskProvider();
    if (!provider) {
      setWalletError('MetaMask not detected. Please install and unlock MetaMask.');
      return;
    }

    try {
      setWalletError('');
      setIsConnecting(true);
      const accounts = await provider.request({ method: 'eth_requestAccounts' });
      setWalletAddress(accounts?.[0] || '');
    } catch (_error) {
      setWalletError('Wallet connection failed or was rejected.');
    } finally {
      setIsConnecting(false);
    }
  };

  const disconnectWallet = async () => {
    const provider = getMetaMaskProvider();
    setWalletError('');
    setIsConnecting(true);

    try {
      if (provider) {
        await provider.request({
          method: 'wallet_revokePermissions',
          params: [{ eth_accounts: {} }],
        });
      }
    } catch (_error) {
      // Ignore revoke failure and still clear local state.
    } finally {
      setWalletAddress('');
      setIsConnecting(false);
    }
  };

  const isContractLab = pathname === '/contract-lab';

  return (
    <div className="app">
      <header className="wallet-bar">
        <div className="wallet-status">{walletAddress ? `Connected: ${shortAddress}` : 'Wallet not connected'}</div>
        <button
          className="wallet-button"
          onClick={walletAddress ? disconnectWallet : connectWallet}
          disabled={isConnecting}
        >
          {isConnecting ? 'Processing...' : walletAddress ? 'Disconnect Wallet' : 'Connect Wallet'}
        </button>
      </header>
      {walletError && <p className="wallet-error">{walletError}</p>}
      {!isContractLab ? (
        <>
          <PricingCards onStartFreeTrial={handleStartFreeTrial} />
        </>
      ) : (
        <ContractLab onBack={() => navigateTo('/', {})} />
      )}
    </div>
  );
}

export default App;