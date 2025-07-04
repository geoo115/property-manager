import axiosInstance from './axiosInstance';

export const getInvoices = async () => {
  const role = localStorage.getItem("role");
  if (role === "admin") {
    const response = await axiosInstance.get("/admin/accounting/invoices");
    return response.data;
  } else if (role === "landlord") {
    const response = await axiosInstance.get("/landlord/invoices");
    return response.data;
  } else if (role === "tenant") {
    const response = await axiosInstance.get("/tenant/invoices");
    return response.data;
  }
  throw new Error("Unauthorized access");
};

export const getInvoiceByID = async (id) => {
  try {
    const response = await axiosInstance.get(`/admin/accounting/invoices/${id}`);
    return response.data;
  } catch (error) {
    console.error("Error fetching invoice by ID:", error);
    throw error;
  }
};

export const createInvoice = async (invoiceData) => {
  const role = localStorage.getItem("role");
  try {
    if (role === "admin") {
      const response = await axiosInstance.post("/admin/accounting/invoices", invoiceData);
      return response.data;
    }
    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error creating invoice:", error);
    throw error;
  }
};

export const updateInvoice = async (id, invoiceData) => {
  const role = localStorage.getItem("role");
  try {
    if (role === "admin") {
      const response = await axiosInstance.put(`/admin/accounting/invoices/${id}`, invoiceData);
      return response.data;
    }
    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error updating invoice:", error);
    throw error;
  }
};

export const deleteInvoice = async (id) => {
  const role = localStorage.getItem("role");
  try {
    if (role === "admin") {
      const response = await axiosInstance.delete(`/admin/accounting/invoices/${id}`);
      return response.data;
    }
    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error deleting invoice:", error);
    throw error;
  }
};

export const getExpenses = async () => {
  const role = localStorage.getItem("role");
  if (role === "admin") {
    const response = await axiosInstance.get("/admin/accounting/expenses");
    return response.data;
  } else if (role === "landlord") {
    const response = await axiosInstance.get("/landlord/expenses");
    return response.data;
  } else if (role === "tenant") {
    throw new Error("Unauthorized: Tenants cannot access expenses directly.");
  }
  throw new Error("Unauthorized access");
};

export const getExpenseByID = async (id) => {
  try {
    const response = await axiosInstance.get(`/admin/accounting/expense/${id}`);
    return response.data;
  } catch (error) {
    console.error("Error fetching expense by ID:", error);
    throw error;
  }
};
export const createExpense = async (expenseData) => {
  const role = localStorage.getItem("role");
  try {
    if (role === "admin" || role === "landlord") {
      return (await axiosInstance.post("/admin/accounting/expense", expenseData)).data;
    }
    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error creating expense:", error);
    throw error;
  }
};

export const updateExpense = async (id, expenseData) => {
  const role = localStorage.getItem("role");
  try {
    if (role === "admin" || role === "landlord") {
      const response = await axiosInstance.put(`/admin/accounting/expense/${id}`, expenseData);
      return response.data;
    }
    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error updating expense:", error);
    throw error;
  }
};

export const deleteExpense = async (id) => {
  const role = localStorage.getItem("role");
  try {
    if (role === "admin") {
      const response = await axiosInstance.delete(`/admin/accounting/expense/${id}`);
      return response.data;
    }
    throw new Error("Unauthorized access");
  } catch (error) {
    console.error("Error deleting expense:", error);
    throw error;
  }
};