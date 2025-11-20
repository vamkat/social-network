'use client'

import { useState } from 'react'

export default function RegisterForm() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    firstName: '',
    lastName: '',
    dateOfBirth: '',
    avatar: '',
    nickname: '',
    aboutMe: '',
  })
  const [errors, setErrors] = useState({})

  const validateForm = () => {
    const newErrors = {}

    // Required field validation
    if (!formData.email.trim()) {
      newErrors.email = 'Email is required'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = 'Please enter a valid email address'
    }

    if (!formData.password.trim()) {
      newErrors.password = 'Password is required'
    } else if (formData.password.length < 8) {
      newErrors.password = 'Password must be at least 8 characters'
    }

    if (!formData.confirmPassword.trim()) {
      newErrors.confirmPassword = 'Please confirm your password'
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match'
    }

    if (!formData.firstName.trim()) {
      newErrors.firstName = 'First name is required'
    }

    if (!formData.lastName.trim()) {
      newErrors.lastName = 'Last name is required'
    }

    if (!formData.dateOfBirth) {
      newErrors.dateOfBirth = 'Date of birth is required'
    } else {
      // Check if user is 13+ years old
      const birthDate = new Date(formData.dateOfBirth)
      const today = new Date()
      const age = today.getFullYear() - birthDate.getFullYear()
      const monthDiff = today.getMonth() - birthDate.getMonth()
      const dayDiff = today.getDate() - birthDate.getDate()

      const actualAge = monthDiff < 0 || (monthDiff === 0 && dayDiff < 0) ? age - 1 : age

      if (actualAge < 13) {
        newErrors.dateOfBirth = 'You must be at least 13 years old to register'
      }
    }

    // Optional field validation
    if (formData.nickname && (formData.nickname.length < 3 || formData.nickname.length > 20)) {
      newErrors.nickname = 'Nickname must be between 3-20 characters'
    }

    if (formData.nickname && !/^[a-zA-Z0-9]+$/.test(formData.nickname)) {
      newErrors.nickname = 'Nickname must be alphanumeric only'
    }

    if (formData.aboutMe && formData.aboutMe.length > 500) {
      newErrors.aboutMe = 'About me must not exceed 500 characters'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const handleSubmit = async (event) => {
    event.preventDefault()

    if (!validateForm()) {
      return
    }

    try {
      const response = await fetch('/api/v1/auth/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email: formData.email,
          password: formData.password,
          firstName: formData.firstName,
          lastName: formData.lastName,
          dateOfBirth: formData.dateOfBirth,
          avatar: formData.avatar || null,
          nickname: formData.nickname || null,
          aboutMe: formData.aboutMe || null,
        }),
      })

      if (!response.ok) {
        const errorData = await response.json()
        if (response.status === 409) {
          setErrors({ email: 'This email is already registered' })
        } else {
          setErrors({ general: errorData.message || 'Registration failed. Please try again.' })
        }
        return
      }

      const data = await response.json()
      console.log(data)
      // Redirect to /feed/public
      window.location.href = '/feed/public'
    } catch (error) {
      console.error('Registration error', error)
      setErrors({ general: 'Registration failed. Please try again.' })
    }
  }

  const handleChange = (event) => {
    const { name, value } = event.target
    setFormData({ ...formData, [name]: value })
    // Clear error when user starts typing
    if (errors[name]) {
      setErrors({ ...errors, [name]: '' })
    }
  }

  const handleImageUpload = (event) => {
    const file = event.target.files[0]
    if (file) {
      // Check file type
      const validTypes = ['image/jpeg', 'image/png', 'image/gif']
      if (!validTypes.includes(file.type)) {
        setErrors({ ...errors, avatar: 'Please upload a JPEG, PNG, or GIF image' })
        return
      }

      // Check file size (max 5MB)
      if (file.size > 5 * 1024 * 1024) {
        setErrors({ ...errors, avatar: 'Image size must be less than 5MB' })
        return
      }

      // Convert to base64
      const reader = new FileReader()
      reader.onloadend = () => {
        setFormData({ ...formData, avatar: reader.result })
        setErrors({ ...errors, avatar: '' })
      }
      reader.readAsDataURL(file)
    }
  }

  const handleRemoveImage = () => {
    setFormData({ ...formData, avatar: '' })
    // Reset the file input
    const fileInput = document.querySelector('input[name="avatar"]')
    if (fileInput) {
      fileInput.value = ''
    }
  }

  return (
    <form className="form-container" onSubmit={handleSubmit}>
      {errors.general && (
        <div className="form-error-general">{errors.general}</div>
      )}

      <div className="form-grid">
        {/* Required Fields */}
        <label className="form-label">
          <span className="form-label-text">
            Email <span className="form-required">*</span>
          </span>
          <input
            type="email"
            name="email"
            value={formData.email}
            placeholder="Enter your email"
            onChange={handleChange}
            className={`form-input ${errors.email ? 'form-input-error' : ''}`}
          />
          {errors.email && <span className="form-error">{errors.email}</span>}
        </label>

        <label className="form-label">
          <span className="form-label-text">
            Password <span className="form-required">*</span>
          </span>
          <input
            type="password"
            name="password"
            value={formData.password}
            placeholder="Min 8 characters"
            onChange={handleChange}
            className={`form-input ${errors.password ? 'form-input-error' : ''}`}
          />
          {errors.password && <span className="form-error">{errors.password}</span>}
        </label>

        <label className="form-label">
          <span className="form-label-text">
            Confirm Password <span className="form-required">*</span>
          </span>
          <input
            type="password"
            name="confirmPassword"
            value={formData.confirmPassword}
            placeholder="Confirm password"
            onChange={handleChange}
            className={`form-input ${errors.confirmPassword ? 'form-input-error' : ''}`}
          />
          {errors.confirmPassword && <span className="form-error">{errors.confirmPassword}</span>}
        </label>

        <label className="form-label">
          <span className="form-label-text">
            First Name <span className="form-required">*</span>
          </span>
          <input
            type="text"
            name="firstName"
            value={formData.firstName}
            placeholder="Enter first name"
            onChange={handleChange}
            className={`form-input ${errors.firstName ? 'form-input-error' : ''}`}
          />
          {errors.firstName && <span className="form-error">{errors.firstName}</span>}
        </label>

        <label className="form-label">
          <span className="form-label-text">
            Last Name <span className="form-required">*</span>
          </span>
          <input
            type="text"
            name="lastName"
            value={formData.lastName}
            placeholder="Enter last name"
            onChange={handleChange}
            className={`form-input ${errors.lastName ? 'form-input-error' : ''}`}
          />
          {errors.lastName && <span className="form-error">{errors.lastName}</span>}
        </label>

        <label className="form-label">
          <span className="form-label-text">
            Date of Birth <span className="form-required">*</span>
          </span>
          <input
            type="date"
            name="dateOfBirth"
            value={formData.dateOfBirth}
            onChange={handleChange}
            className={`form-input form-input-date ${errors.dateOfBirth ? 'form-input-error' : ''}`}
          />
          {errors.dateOfBirth && <span className="form-error">{errors.dateOfBirth}</span>}
        </label>

        {/* Optional Fields */}
        <div className="form-label">
          <span className="form-label-text">Avatar/Image (Optional)</span>
          <input
            type="file"
            name="avatar"
            accept="image/jpeg,image/png,image/gif"
            onChange={handleImageUpload}
            className="form-input-file"
          />
          {formData.avatar && (
            <div className="form-image-preview">
              <div className="relative inline-block">
                <img src={formData.avatar} alt="Avatar preview" className="form-avatar-preview" />
                <button
                  type="button"
                  onClick={handleRemoveImage}
                  className="absolute -top-2 -right-2 w-6 h-6 bg-red-500/80 hover:bg-red-500 text-white rounded-full flex items-center justify-center text-sm font-bold transition-all duration-200 border border-white/20"
                  aria-label="Remove image"
                >
                  ×
                </button>
              </div>
            </div>
          )}
          {errors.avatar && <span className="form-error">{errors.avatar}</span>}
        </div>

        <label className="form-label">
          <span className="form-label-text">Nickname (Optional)</span>
          <input
            type="text"
            name="nickname"
            value={formData.nickname}
            placeholder="3-20 alphanumeric characters"
            onChange={handleChange}
            className={`form-input ${errors.nickname ? 'form-input-error' : ''}`}
          />
          {errors.nickname && <span className="form-error">{errors.nickname}</span>}
        </label>

        <label className="form-label form-label-span-2">
          <span className="form-label-text">About Me (Optional)</span>
          <textarea
            name="aboutMe"
            value={formData.aboutMe}
            placeholder="Tell us about yourself (max 500 characters)"
            rows="4"
            onChange={handleChange}
            className={`form-textarea ${errors.aboutMe ? 'form-input-error' : ''}`}
          />
          <span className="form-char-count">{formData.aboutMe.length}/500</span>
          {errors.aboutMe && <span className="form-error">{errors.aboutMe}</span>}
        </label>
      </div>

      <div className="hero-cta-group flex justify-center-safe items-center p-6">
        <button className="btn-primary" type="submit">
          <span>Create Account</span>
        </button>

        <button className="btn-secondary" type="button" onClick={() => window.location.href = '/login'}>
          Cancel
        </button>
      </div>
    </form>
  )
}
