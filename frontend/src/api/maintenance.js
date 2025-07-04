import axiosInstance from '../api/axiosInstance';

export const getMaintenances = async (leaseId = null, propertyId = null) => {
  const role = localStorage.getItem("role");
  
  try {
    let response;
    switch(role) {
      case "admin":
        response = await axiosInstance.get("/admin/maintenances");
        return response.data;
      case "maintenanceTeam":
        response = await axiosInstance.get("/maintenanceTeam/maintenances");
        return response.data;
      case "tenant":
        if (!leaseId) throw new Error("Missing lease ID");
        response = await axiosInstance.get(`/tenant/leases/${leaseId}/maintenance`);
        return response.data; 
      case "landlord":
        if (!propertyId) throw new Error("Missing property ID");
        response = await axiosInstance.get(`/landlord/properties/${propertyId}/maintenances`);
        return response.data;
      default:
        throw new Error("Unauthorized access");
    }
  } catch (error) {
    console.error("Error fetching maintenances:", error);
    throw error;
  }
};

export const createMaintenance = async (leaseId, maintenanceData, propertyId = null) => {
  const role = localStorage.getItem("role");
  
  try {
    let response;
    switch(role) {
      case 'tenant':
        if (!leaseId) throw new Error("Missing lease ID");
        response = await axiosInstance.post(`/tenant/leases/${leaseId}/maintenance`, maintenanceData);
        break;
      case 'admin':
        if (!propertyId) throw new Error("Missing property ID");
        response = await axiosInstance.post(`/admin/properties/${propertyId}/maintenances`, maintenanceData);
        break;
      case 'landlord':
        if (!propertyId) throw new Error("Missing property ID");
        response = await axiosInstance.post(`/landlord/properties/${propertyId}/maintenances`, maintenanceData);
        break;
      default:
        throw new Error("Unauthorized access");
    }
    return response.data;
  } catch (error) {
    console.error("Error creating maintenance:", error);
    throw error;
  }
};

export const updateMaintenance = async (id, maintenanceData) => {
  const role = localStorage.getItem("role");

  try {
    if (role === "maintenanceTeam") {
      return (await axiosInstance.put(`/maintenanceTeam/maintenance/${id}`, maintenanceData)).data;
    }
    if (role === "admin") {
      return (await axiosInstance.put(`/admin/maintenances/${id}`, maintenanceData)).data;
    }

    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error updating maintenance:", error);
    throw error;
  }
};

export const deleteMaintenance = async (id) => {
  const role = localStorage.getItem("role");

  try {
    if (role !== "admin") throw new Error("Unauthorized");

    return (await axiosInstance.delete(`/admin/maintenances/${id}`)).data;
  } catch (error) {
    console.error("Error deleting maintenance:", error);
    throw error;
  }
};