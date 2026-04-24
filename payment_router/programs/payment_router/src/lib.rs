use anchor_lang::prelude::*;
use anchor_lang::system_program::{transfer as system_transfer, Transfer as SystemTransfer};
use anchor_spl::token::{self, Mint, Token, TokenAccount, TransferChecked};

declare_id!("Cv1XLdNwF2Hk6Ldj2U2c94dJ4i29vHfTRFs1vzKqDDMs");

#[program]
pub mod payment_router {
    use super::*;

    /// Move native SOL from `payer` to `recipient`. The transaction builder supplies `recipient`.
    pub fn transfer_sol(ctx: Context<TransferSol>, lamports: u64) -> Result<()> {
        require!(lamports > 0, PaymentRouterError::ZeroAmount);

        let cpi = SystemTransfer {
            from: ctx.accounts.payer.to_account_info(),
            to: ctx.accounts.recipient.to_account_info(),
        };
        system_transfer(
            CpiContext::new(ctx.accounts.system_program.to_account_info(), cpi),
            lamports,
        )?;
        Ok(())
    }

    /// Move SPL tokens from payer's ATA to recipient's ATA. Builder passes both ATAs + mint.
    pub fn transfer_spl(ctx: Context<TransferSpl>, amount: u64) -> Result<()> {
        require!(amount > 0, PaymentRouterError::ZeroAmount);

        let decimals = ctx.accounts.mint.decimals;

        let cpi = TransferChecked {
            from: ctx.accounts.from.to_account_info(),
            mint: ctx.accounts.mint.to_account_info(),
            to: ctx.accounts.to.to_account_info(),
            authority: ctx.accounts.authority.to_account_info(),
        };
        token::transfer_checked(
            CpiContext::new(ctx.accounts.token_program.to_account_info(), cpi),
            amount,
            decimals,
        )?;
        Ok(())
    }
}

#[derive(Accounts)]
pub struct TransferSol<'info> {
    #[account(mut)]
    pub payer: Signer<'info>,
    /// Recipient wallet; not a PDA — whoever builds the tx sets this account.
    /// CHECK: system-owned lamport destination
    #[account(mut)]
    pub recipient: UncheckedAccount<'info>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct TransferSpl<'info> {
    #[account(mut)]
    pub mint: Account<'info, Mint>,
    #[account(
        mut,
        constraint = from.owner == authority.key() @ PaymentRouterError::InvalidFromOwner,
        constraint = from.mint == mint.key() @ PaymentRouterError::MintMismatch
    )]
    pub from: Account<'info, TokenAccount>,
    #[account(mut, constraint = to.mint == mint.key() @ PaymentRouterError::MintMismatch)]
    pub to: Account<'info, TokenAccount>,
    pub authority: Signer<'info>,
    pub token_program: Program<'info, Token>,
}

#[error_code]
pub enum PaymentRouterError {
    #[msg("amount must be > 0")]
    ZeroAmount,
    #[msg("from token account owner must match authority")]
    InvalidFromOwner,
    #[msg("from/to mint must match mint account")]
    MintMismatch,
}
