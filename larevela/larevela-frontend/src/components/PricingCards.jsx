// PricingCards.jsx
import React from 'react';
import './PricingCards.css';

const PricingCards = () => {
  const plans = [
    {
      name: "Growth",
      price: "$199",
      period: "/month",
      description: "For growing businesses",
      extraInfo: "+ $0.05 per identity revealed",
      type: "growth",
      features: [
        "Identity Enrichment",
        "Visitor Analytics",
        "Audience Builder",
        "Bot Identification",
        "Traffic Validation",
        "Email Verification",
        "Email Export (CSV)"
      ],
      buttonText: "Start Free Trial",
      buttonType: "outline",
      highlight: true
    },
    {
      name: "Professional",
      price: "$399",
      period: "/month",
      description: "For established businesses",
      extraInfo: "+ $0.05 per identity revealed",
      type: "professional",
      features: [
        "Everything in Growth",
        { divider: true, label: "Plus" },
        "Customer Journey Events",
        "Sync Sources",
        "Integrations Auto Sync",
        "Multi-User",
        "Multiple Websites",
        "Reporting"
      ],
      buttonText: "Start Free Trial",
      buttonType: "outline",
      highlight: false
    },
    {
      name: "Enterprise",
      price: "Custom",
      period: "Contact Sales",
      description: "For large organizations",
      extraInfo: "+ $0.05 per identity revealed",
      type: "enterprise",
      features: [
        "Everything in Professional",
        { divider: true, label: "Plus" },
        "Reporting Automation",
        "AI Audience Builder",
        "Single Sign On",
        "Custom contract terms",
        "Dedicated account manager",
        "SLA guarantee"
      ],
      buttonText: "Let's Talk",
      buttonType: "solid",
      highlight: false
    }
  ];

  const CheckIcon = () => (
    <svg className="icon-check" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 6 9 17l-5-5"></path>
    </svg>
  );

  const ShieldIcon = () => (
    <svg className="icon-shield" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M20 13c0 5-3.5 7.5-7.66 8.95a1 1 0 0 1-.67-.01C7.5 20.5 4 18 4 13V6a1 1 0 0 1 1-1c2 0 4.5-1.2 6.24-2.72a1.17 1.17 0 0 1 1.52 0C14.51 3.81 17 5 19 5a1 1 0 0 1 1 1z"></path>
    </svg>
  );

  return (
    <section className="pricing-section">
      <div className="pricing-container">
        <div className="pricing-grid">
          {plans.map((plan, index) => (
            <div key={index} className="pricing-card-wrapper">
              <div className={`pricing-card ${plan.highlight ? 'highlight' : 'normal'}`}>
                {/* Card Header */}
                <div className={`card-header ${plan.highlight ? 'highlight' : ''}`}>
                  <h4 className="card-title">{plan.name}</h4>
                  <div className="price-wrapper">
                    <div className={`price-box ${plan.type}`}>
                      <div className="price">{plan.price}</div>
                      <div className="price-period">{plan.period}</div>
                    </div>
                  </div>
                  <p className="card-description">{plan.description}</p>
                  <div className="extra-info">
                    <div className="extra-info-text">{plan.extraInfo}</div>
                  </div>
                </div>

                {/* Card Content */}
                <div className="card-content">
                  <ul className="features-list">
                    {plan.features.map((feature, idx) => {
                      if (typeof feature === 'object' && feature.divider) {
                        return (
                          <li key={idx}>
                            <div className="divider-wrapper">
                              <div className="divider-line"></div>
                              <span className="divider-label">{feature.label}</span>
                              <div className="divider-line"></div>
                            </div>
                          </li>
                        );
                      }
                      return (
                        <li key={idx}>
                          <CheckIcon />
                          <span className="feature-text">{feature}</span>
                        </li>
                      );
                    })}
                  </ul>

                  <div className="button-wrapper">
                    {plan.buttonType === 'outline' && (
                      <>
                        <div className="shield-note">
                          <ShieldIcon />
                          No credit card required
                        </div>
                        <button className="btn btn-outline">
                          {plan.buttonText}
                        </button>
                      </>
                    )}
                    {plan.buttonType === 'solid' && (
                      <button className="btn btn-solid">
                        {plan.buttonText}
                      </button>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default PricingCards;