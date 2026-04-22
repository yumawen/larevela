import React, { useEffect, useMemo, useRef, useState } from 'react';
import {
  Connection,
  LAMPORTS_PER_SOL,
  PublicKey,
  clusterApiUrl,
} from '@solana/web3.js';
import './ContractLab.css';

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
  const [readResult, setReadResult] = useState('');
  const [readLoading, setReadLoading] = useState(false);
  const [paymentNo, setPaymentNo] = useState('');
  const [paymentStatusLoading, setPaymentStatusLoading] = useState(false);
  const [paymentStatusError, setPaymentStatusError] = useState('');
  const [paymentStatus, setPaymentStatus] = useState(null);
  const [lastRefreshedAt, setLastRefreshedAt] = useState('');
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
    setReadResult('');

    try {
      const { address } = await ensureWalletConnected();
      const lamports = await connection.getBalance(new PublicKey(address));
      const balanceSol = (lamports / LAMPORTS_PER_SOL).toFixed(6);
      setReadResult(`Wallet balance: ${balanceSol} SOL (devnet) | Address: ${address}`);
    } catch (error) {
      setReadResult(error?.message || 'Query failed. Please retry.');
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

  return (
    <section className="contract-lab">
      <div className="lab-header">
        <button className="lab-back" onClick={onBack}>
          Back to Pricing
        </button>
        <h2>Contract Interaction Lab (Solflare Only)</h2>
      </div>

      <div className="lab-card">
        <h3>1) Query wallet balance</h3>
        <p className="lab-muted">
          Read-only devnet RPC: SOL balance for the connected Solflare wallet address.
        </p>
        <button className="lab-btn" onClick={queryWalletBalance} disabled={readLoading}>
          {readLoading ? 'Querying...' : 'Query balance'}
        </button>
        {readResult && <p className="lab-result">{readResult}</p>}
      </div>

      <div className="lab-card">
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
    </section>
  );
};

export default ContractLab;
