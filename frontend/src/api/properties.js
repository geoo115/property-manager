import axiosInstance from '../api/axiosInstance';

export const getProperties = async () => {
  const role = localStorage.getItem("role");

  if (role === "admin") {
    return (await axiosInstance.get("/admin/properties")).data;
  } else if (role === "landlord") {
    return (await axiosInstance.get("/landlord/properties")).data;
  } else if (role === "tenant") {
    throw new Error("Unauthorized: Tenants cannot access properties directly.");
  }

  throw new Error("Unauthorized access");
};
;

export const getPropertyByID = async (id) => {
  const role = localStorage.getItem("role");
  if (role === "admin") {
    return (await axiosInstance.get(`/admin/properties/${id}`)).data;
  } else if (role === "landlord") {
    return (await axiosInstance.get(`/landlord/properties/${id}`)).data;
  }
  throw new Error("Unauthorized access");
};

// Only Admin can create/update/delete properties
export const createProperty = async (propertyData) => {
  enforceAdmin();
  return (await axiosInstance.post("/admin/properties", propertyData)).data;
};

export const updateProperty = async (id, propertyData) => {
  enforceAdmin();
  return (await axiosInstance.put(`/admin/properties/${id}`, propertyData)).data;
};

export const deleteProperty = async (id) => {
  enforceAdmin();
  return (await axiosInstance.delete(`/admin/properties/${id}`)).data;
};

// Utility function for role enforcement
const enforceAdmin = () => {
  if (localStorage.getItem("role") !== "admin") throw new Error("Unauthorized");
};
