import { useState } from 'react';
import axios from 'axios';

export const useSearchKeyword = () => {
    const [filepath, setFilepath] = useState(null);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);

    const searchKeyword = async (keyword) => {
        try {
            setLoading(true);
            const response = await axios.get(`${import.meta.env.VITE_DOMAINS_LAB_API}/search`, {
                params: { keyword },
            });
            setFilepath(response.data.filepath);
        } catch (error) {
            setError(error.message);
            console.log(error);
        } finally {
            setLoading(false);
        }
    };

    const downloadFile = async (filePath) => {
        try {
          const response = await axios.get(`${import.meta.env.VITE_DOMAINS_LAB_API}/download`, {
            params: { file: filePath },
            responseType: 'blob',
          });
          const blob = new Blob([response.data], { type: 'application/octet-stream' });
          const url = URL.createObjectURL(blob);
          const a = document.createElement('a');
          a.href = url;
          a.download = filePath;
          a.click();
        } catch (error) {
          console.error(error);
        }
      };

    return { filepath, error, searchKeyword, downloadFile, loading };
};