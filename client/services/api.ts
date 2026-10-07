import axios from '@/lib/axios';

// Types
export interface UploadFileResponse {
  success: boolean;
  fileUrl: string;
  fileName: string;
  message?: string;
}

export interface FileInfo {
  id: string;
  name: string;
  url: string;
  size: number;
  uploadedAt: string;
}

// API Service
export const apiService = {
  // Upload file to R2
  uploadFile: async (file: File): Promise<UploadFileResponse> => {
    const formData = new FormData();
    formData.append('file', file);

    const response = await axios.post<UploadFileResponse>('/api/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Get list of files
  getFiles: async (): Promise<FileInfo[]> => {
    const response = await axios.get<FileInfo[]>('/api/files');
    return response.data;
  },

  // Get file by ID
  getFile: async (fileId: string): Promise<FileInfo> => {
    const response = await axios.get<FileInfo>(`/api/files/${fileId}`);
    return response.data;
  },

  // Delete file
  deleteFile: async (fileId: string): Promise<void> => {
    await axios.delete(`/api/files/${fileId}`);
  },

  // Get presigned URL for upload
  getPresignedUrl: async (fileName: string): Promise<{ url: string; key: string }> => {
    const response = await axios.post<{ url: string; key: string }>('/api/presigned-url', {
      fileName,
    });
    return response.data;
  },

  // Health check
  healthCheck: async (): Promise<{ status: string }> => {
    const response = await axios.get<{ status: string }>('/api/health');
    return response.data;
  },
};

export default apiService;
