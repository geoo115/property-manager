import React, { useContext, useState, useEffect } from 'react';
import { AuthContext } from '../context/AuthContext';
import { TenantOnly } from '../components/common/RoleBasedContent';
import Button from '../components/common/Button';
import { toast } from 'react-toastify';
import axiosInstance from '../api/axiosInstance';

const TenantProfile = () => {
  const { user, hasUserPermission } = useContext(AuthContext);
  const [profile, setProfile] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState({});

  useEffect(() => {
    fetchProfile();
  }, []);

  const fetchProfile = async () => {
    try {
      const response = await axiosInstance.get('/tenant/profile');
      setProfile(response.data.profile);
      setFormData(response.data.profile || {});
    } catch (error) {
      toast.error('Failed to load profile');
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!hasUserPermission('profile', 'update')) {
      toast.error('You do not have permission to update your profile');
      return;
    }

    try {
      await axiosInstance.put('/tenant/profile', formData);
      setProfile(formData);
      setIsEditing(false);
      toast.success('Profile updated successfully');
    } catch (error) {
      toast.error('Failed to update profile');
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  if (isLoading) {
    return (
      <div className="page-container">
        <div className="loading-state">Loading profile...</div>
      </div>
    );
  }

  return (
    <TenantOnly>
      <div className="page-container">
        <div className="page-header">
          <div>
            <h1>My Profile</h1>
            <p className="page-description">View and update your personal information</p>
          </div>
          {!isEditing && (
            <Button
              variant="primary"
              onClick={() => setIsEditing(true)}
            >
              Edit Profile
            </Button>
          )}
        </div>

        <div className="profile-content">
          {isEditing ? (
            <form onSubmit={handleSubmit} className="profile-form">
              <div className="form-group">
                <label htmlFor="first_name">First Name</label>
                <input
                  type="text"
                  id="first_name"
                  name="first_name"
                  value={formData.first_name || ''}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="last_name">Last Name</label>
                <input
                  type="text"
                  id="last_name"
                  name="last_name"
                  value={formData.last_name || ''}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="email">Email</label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={formData.email || ''}
                  onChange={handleInputChange}
                  required
                />
              </div>

              <div className="form-group">
                <label htmlFor="phone">Phone</label>
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  value={formData.phone || ''}
                  onChange={handleInputChange}
                />
              </div>

              <div className="form-group">
                <label htmlFor="emergency_contact">Emergency Contact</label>
                <input
                  type="text"
                  id="emergency_contact"
                  name="emergency_contact"
                  value={formData.emergency_contact || ''}
                  onChange={handleInputChange}
                />
              </div>

              <div className="form-actions">
                <Button type="submit" variant="primary">
                  Save Changes
                </Button>
                <Button 
                  type="button" 
                  variant="secondary" 
                  onClick={() => {
                    setIsEditing(false);
                    setFormData(profile || {});
                  }}
                >
                  Cancel
                </Button>
              </div>
            </form>
          ) : (
            <div className="profile-display">
              <div className="profile-section">
                <h3>Personal Information</h3>
                <div className="profile-field">
                  <span className="field-label">Name:</span>
                  <span className="field-value">
                    {profile?.first_name} {profile?.last_name}
                  </span>
                </div>
                <div className="profile-field">
                  <span className="field-label">Email:</span>
                  <span className="field-value">{profile?.email}</span>
                </div>
                <div className="profile-field">
                  <span className="field-label">Phone:</span>
                  <span className="field-value">{profile?.phone || 'Not provided'}</span>
                </div>
                <div className="profile-field">
                  <span className="field-label">Emergency Contact:</span>
                  <span className="field-value">{profile?.emergency_contact || 'Not provided'}</span>
                </div>
              </div>

              <div className="profile-section">
                <h3>Account Information</h3>
                <div className="profile-field">
                  <span className="field-label">Role:</span>
                  <span className="field-value role-badge role-tenant">Tenant</span>
                </div>
                <div className="profile-field">
                  <span className="field-label">Member Since:</span>
                  <span className="field-value">
                    {profile?.created_at ? new Date(profile.created_at).toLocaleDateString() : 'Unknown'}
                  </span>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </TenantOnly>
  );
};

export default TenantProfile;
