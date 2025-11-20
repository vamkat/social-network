import RegisterForm from '@/components/features/RegisterForm'

const RegisterPage = () => {

  return (
    <div className="landing-container">
      {/* Grain texture overlay */}
      <div className="grain-texture"></div>
      <section className="hero-section">
        <div className="hero-grid-full">
          {/* Main glass surface - Full Width */}
          <div className="glass-card-wrapper">
            <div className="glass-glow"></div>
            <div className="glass-card-main">


              <h1 className="hero-heading text-center">
                Register to Connect<span className="hero-gradient-text">.</span><br />
              </h1>

              <RegisterForm />

              <div className="hero-cta-group mt-6 flex justify-center gap-4 py-2">
                <span className="text-white/80 text-[15px] font-light">
                  Already registered? <a href="/login" className="border-b border-white/30 hover:border-white/60 pb-0.5 transition-all duration-200">Login Now</a>
                </span>
              </div>

            </div>
            
            {/* Floating accent line */}
            <div className="floating-accent-line"></div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default RegisterPage
