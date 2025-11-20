'use client'

import { useState } from 'react'

export default function LoginForm() {
  const [formData, setFormData] = useState({
    identifier: '',
    password: '',
  })
  const [error, setError] = useState('')

  const handleSubmit = async (event) => {
    event.preventDefault()
    setError('')

    // Client-side validation
    if (!formData.identifier || !formData.password) {
      setError('Please enter both email/username and password')
      return
    }

    try {
      const response = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      })

      if (!response.ok) {
        const errorData = await response.json()

        if (response.status === 404) {
          setError('No account found with this email/username')
        } else if (response.status === 401) {
          setError('Invalid email or password')
        } else {
          setError(errorData.message || 'Login failed. Please try again.')
        }
        return
      }

      const data = await response.json()
      console.log(data)
      // Redirect to /feed/public
      window.location.href = '/feed/public'
    } catch (error) {
      console.error('Login error', error)
      setError('Login failed. Please try again.')
    }
  }

  const handleChange = (event) => {
    const { name, value } = event.target
    setFormData({ ...formData, [name]: value })
    // Clear error when user starts typing
    if (error) {
      setError('')
    }
  }

  return (
    <form className="form-container-center" onSubmit={handleSubmit}>
      {error && (
        <div className="form-error-general">{error}</div>
      )}

      <input
        type="text"
        name="identifier"
        value={formData.identifier}
        placeholder="Email/Username"
        onChange={handleChange}
        className="form-input-center"
        required
      />
      <br />
      <input
        type="password"
        name="password"
        value={formData.password}
        placeholder="Password"
        onChange={handleChange}
        className="form-input-center"
        required
      />

      <div className="hero-cta-group flex justify-center-safe items-center p-6">
        <button className="btn-primary" type="submit">
          <span>Login</span>
        </button>

        <button className="btn-secondary" type="button" onClick={() => window.location.href = '/'}>
          Cancel
        </button>
      </div>
    </form>
  )
}
