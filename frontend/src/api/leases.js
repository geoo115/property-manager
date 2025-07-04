import axiosInstance from '../api/axiosInstance';

// Fetch leases based on role
export const getLeases = async (params) => {
  const role = localStorage.getItem("role");

  if (role === "admin") {
    return (await axiosInstance.get("/admin/leases", { params })).data;
  } else if (role === "landlord") {
    return (await axiosInstance.get("/landlord/leases", { params })).data;
  } else if (role === "tenant") {
    return (await axiosInstance.get("/tenant/leases", { params })).data;
  }

  throw new Error("Unauthorized");
};

export const getLeaseByID = async (id) => {
  const role = localStorage.getItem("role");
  if (role === "admin") {
    return (await axiosInstance.get(`/admin/leases/${id}`)).data;
  } else if (role === "landlord") {
    return (await axiosInstance.get(`/landlord/leases/${id}`)).data;
  } else if (role === "tenant") {
    return (await axiosInstance.get(`/tenant/leases/${id}`)).data;
  }
  throw new Error("Unauthorized");
};

// Only Admin can create/update/delete leases
export const createLease = async (leaseData) => {
  enforceAdmin();
  return (await axiosInstance.post("/admin/leases", leaseData)).data;
};

export const updateLease = async (id, leaseData) => {
  enforceAdmin();
  return (await axiosInstance.put(`/admin/leases/${id}`, leaseData)).data;
};

export const deleteLease = async (id) => {
  enforceAdmin();
  return (await axiosInstance.delete(`/admin/leases/${id}`)).data;
};

// Utility function for role enforcement
const enforceAdmin = () => {
  if (localStorage.getItem("role") !== "admin") {
    throw new Error("Unauthorized");
  }
};
