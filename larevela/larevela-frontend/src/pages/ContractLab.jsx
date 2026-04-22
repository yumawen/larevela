import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  Connection,
  LAMPORTS_PER_SOL,
  PublicKey,
  clusterApiUrl,
} from '@solana/web3.js';
import './ContractLab.css';

const DEVNET_USDC_MINT = '4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU';
const TRIAL_USDC_RECEIVER = '7xRr4GRzw5aw441Btum3Zxot6RUVBEUGihMdkwFb17zc';
const TOKEN_PROGRAM_PUBKEY = new PublicKey('TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA');
const ASSOCIATED_TOKEN_PROGRAM_PUBKEY = new PublicKey('ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL');
const truncateAddress = (address) => {
  if (!address) {
    return '-';
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

const ContractLab = ({ onBack }) => {
  const connection = useMemo(() => new Connection(clusterApiUrl('devnet'), 'confirmed'), []);
  const [readResult, setReadResult] = useState(null);
  const [readError, setReadError] = useState('');
  const [readLoading, setReadLoading] = useState(false);
  const [paymentNo, setPaymentNo] = useState('');
  const [paymentStatusLoading, setPaymentStatusLoading] = useState(false);
  const [paymentStatusError, setPaymentStatusError] = useState('');
  const [paymentStatus, setPaymentStatus] = useState(null);
  const [lastRefreshedAt, setLastRefreshedAt] = useState('');
  const [contractReadLoading, setContractReadLoading] = useState(false);
  const [contractReadResult, setContractReadResult] = useState(null);
  const [contractReadError, setContractReadError] = useState('');
  const paymentRefreshInFlightRef = useRef(false);

  useEffect(() => {
    const savedPaymentNo = localStorage.getItem('lastPaymentNo') || '';
    if (savedPaymentNo) {
      setPaymentNo(savedPaymentNo);
    }
  }, []);

  const ensureWalletConnected = async () => {
    const provider = getSolflareProvider();
    if (!provider) {
      throw new Error('Solflare Wallet not found. Please install and unlock Solflare.');
    }

    if (provider.publicKey) {
      return { provider, address: provider.publicKey.toBase58() };
    }

    const response = await provider.connect();
    const address = response?.publicKey?.toBase58?.() || provider.publicKey?.toBase58?.();
    if (!address) {
      throw new Error('No account returned by Solflare.');
    }
    return { provider, address };
  };

  const queryWalletBalance = async () => {
    setReadLoading(true);
    setReadResult(null);
    setReadError('');

    try {
      const { address } = await ensureWalletConnected();
      const owner = new PublicKey(address);
      const lamports = await connection.getBalance(owner);
      const balanceSol = (lamports / LAMPORTS_PER_SOL).toFixed(6);

      const tokenResp = await connection.getParsedTokenAccountsByOwner(
        owner,
        { programId: TOKEN_PROGRAM_PUBKEY },
        'confirmed'
      );
      const tokenBalances = tokenResp.value
        .map((item) => {
          const info = item.account.data?.parsed?.info || {};
          const tokenAmount = info.tokenAmount?.uiAmountString ?? info.tokenAmount?.amount ?? '0';
          return {
            mint: info.mint || '',
            amount: tokenAmount,
            ata: item.pubkey.toBase58(),
          };
        })
        .filter((item) => item.mint)
        .slice(0, 8);

      setReadResult({
        address,
        balanceSol,
        tokenBalances,
      });
    } catch (error) {
      setReadError(error?.message || 'Query failed. Please retry.');
    } finally {
      setReadLoading(false);
    }
  };

  const refreshPaymentStatus = async ({ silent = false } = {}) => {
    const normalizedPaymentNo = paymentNo.trim();
    if (!normalizedPaymentNo) {
      setPaymentStatus(null);
      setPaymentStatusError('Please input payment number first.');
      return;
    }

    if (paymentRefreshInFlightRef.current) {
      return;
    }
    paymentRefreshInFlightRef.current = true;
    if (!silent) {
      setPaymentStatusLoading(true);
    }
    setPaymentStatusError('');

    try {
      const resp = await fetch(`http://localhost:8888/api/v1/trade/payments/${encodeURIComponent(normalizedPaymentNo)}`);
      if (!resp.ok) {
        throw new Error('Failed to query payment status.');
      }
      const data = await resp.json();
      setPaymentStatus(data);
      localStorage.setItem('lastPaymentNo', normalizedPaymentNo);
      setLastRefreshedAt(new Date().toLocaleTimeString());
    } catch (error) {
      setPaymentStatus(null);
      setPaymentStatusError(error?.message || 'Failed to refresh payment status.');
    } finally {
      paymentRefreshInFlightRef.current = false;
      if (!silent) {
        setPaymentStatusLoading(false);
      }
    }
  };

  const queryReceiverUsdcAtaInfo = async () => {
    setContractReadLoading(true);
    setContractReadResult(null);
    setContractReadError('');
    try {
      const mintPubkey = new PublicKey(DEVNET_USDC_MINT);
      const receiverPubkey = new PublicKey(TRIAL_USDC_RECEIVER);
      const [receiverAta] = PublicKey.findProgramAddressSync(
        [receiverPubkey.toBuffer(), TOKEN_PROGRAM_PUBKEY.toBuffer(), mintPubkey.toBuffer()],
        ASSOCIATED_TOKEN_PROGRAM_PUBKEY
      );

      const accountInfo = await connection.getParsedAccountInfo(receiverAta, 'confirmed');
      if (!accountInfo.value) {
        setContractReadError(`Receiver USDC ATA not initialized yet. ATA: ${receiverAta.toBase58()}`);
        return;
      }

      const parsed = accountInfo.value.data?.parsed?.info || {};
      const tokenAmount = parsed.tokenAmount?.uiAmountString || parsed.tokenAmount?.amount || '0';
      setContractReadResult({
        ata: receiverAta.toBase58(),
        owner: parsed.owner || '-',
        mint: parsed.mint || '-',
        amount: tokenAmount,
      });
    } catch (error) {
      setContractReadError(error?.message || 'Failed to read contract account info.');
    } finally {
      setContractReadLoading(false);
    }
  };

  return (
    <section className="contract-lab">
      <div className="lab-header">
        <button className="lab-back" onClick={onBack}>
          Back to Pricing
        </button>
        <h2>Contract Interaction Lab (Solflare Only)</h2>
      </div>

      <div className="lab-card lab-card--wallet">
        <h3>1) Query wallet balance</h3>
        <p className="lab-muted">
          Read-only devnet RPC: SOL and SPL token balances for the connected Solflare wallet.
        </p>
        <button className="lab-btn" onClick={queryWalletBalance} disabled={readLoading}>
          {readLoading ? 'Querying...' : 'Query balance'}
        </button>
        {readError && <p className="lab-result">{readError}</p>}
        {readResult && (
          <div className="lab-result">
            <p>Address: <span title={readResult.address}>{truncateAddress(readResult.address)}</span></p>
            <p>SOL balance: {readResult.balanceSol} SOL</p>
            <p>SPL token balances (up to 8 accounts):</p>
            <div className="lab-list">
              {readResult.tokenBalances.length === 0 ? (
                <p>No SPL token account found.</p>
              ) : (
                readResult.tokenBalances.map((item) => (
                  <p key={`${item.ata}-${item.mint}`}>
                    Mint: <span title={item.mint}>{truncateAddress(item.mint)}</span> | Amount: {item.amount} | ATA: <span title={item.ata}>{truncateAddress(item.ata)}</span>
                  </p>
                ))
              )}
            </div>
          </div>
        )}
      </div>

      <div className="lab-card lab-card--payment">
        <h3>2) Query payment status (backend)</h3>
        <p className="lab-muted">Refresh payment status from trade-api aggregated view.</p>
        <label className="lab-label" htmlFor="paymentNo">
          Payment number
        </label>
        <input
          id="paymentNo"
          className="lab-input"
          placeholder="pay-..."
          value={paymentNo}
          onChange={(event) => setPaymentNo(event.target.value)}
        />
        <button className="lab-btn" onClick={refreshPaymentStatus} disabled={paymentStatusLoading}>
          {paymentStatusLoading ? 'Refreshing...' : 'Refresh payment status'}
        </button>
        {lastRefreshedAt && <p className="lab-muted">Last refreshed at: {lastRefreshedAt}</p>}
        {paymentStatusError && <p className="lab-result">{paymentStatusError}</p>}
        {paymentStatus && (
          <div className="lab-result">
            <p>PaymentNo: {paymentStatus.paymentNo || '-'}</p>
            <p>OrderNo: {paymentStatus.orderNo || '-'}</p>
            <p>Status: {paymentStatus.status || '-'}</p>
            <p>Confirmation: {paymentStatus.confirmationStatus || '-'} ({paymentStatus.confirmations ?? 0})</p>
            <p>TxId: {paymentStatus.txId || '-'}</p>
            <p>Amount: {paymentStatus.amountActual || '-'} / expected {paymentStatus.amountExpected || '-'}</p>
            <p>From: {paymentStatus.payerAccount || '-'}</p>
            <p>To: {paymentStatus.receiverAccount || '-'}</p>
            {paymentStatus.failureReason && <p>Failure reason: {paymentStatus.failureReason}</p>}
          </div>
        )}
      </div>

      <div className="lab-card lab-card--contract">
        <h3>3) Query contract account info (SPL Token Account)</h3>
        <p className="lab-muted">
          Read on-chain receiver USDC associated token account state from Solana devnet.
        </p>
        <button className="lab-btn" onClick={queryReceiverUsdcAtaInfo} disabled={contractReadLoading}>
          {contractReadLoading ? 'Querying...' : 'Query receiver USDC ATA'}
        </button>
        {contractReadError && <p className="lab-result">{contractReadError}</p>}
        {contractReadResult && (
          <div className="lab-result">
            <p>Receiver ATA: <span title={contractReadResult.ata}>{truncateAddress(contractReadResult.ata)}</span></p>
            <p>Owner: <span title={contractReadResult.owner}>{truncateAddress(contractReadResult.owner)}</span></p>
            <p>Mint: <span title={contractReadResult.mint}>{truncateAddress(contractReadResult.mint)}</span></p>
            <p>USDC balance: {contractReadResult.amount}</p>
          </div>
        )}
      </div>
    </section>
  );
};

export default ContractLab;
