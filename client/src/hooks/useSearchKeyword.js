import { useState } from 'react';
import axios from 'axios';

export const useSearchKeyword = () => {
    const [filepath, setFilepath] = useState(null);
    const [error, setError] = useState(null);
    const [loading, setLoading] = useState(false);

    const searchKeyword = async (keyword) => {
        try {
            setLoading(true);
            setError(null);

            // Make the API call to search the database
            const response = await axios.get(`${import.meta.env.VITE_DOMAINS_LAB_API}/search`, {
                params: { keyword },
            });
            // Set the filepath to the result from the API
            setFilepath(response.data.filepath);
        } catch (error) {
            setError(error.response.data);
            console.log(error);
        } finally {
            setLoading(false);
        }
    };


    const downloadFile = async (filePath) => {
        try {
          // Make the API call to download the file
          const response = await axios.get(`${import.meta.env.VITE_DOMAINS_LAB_API}/download`, {
            params: { file: filePath },
            responseType: 'blob',
          });

          // Create a blob to download the file from the response data
          const blob = new Blob([response.data], { type: 'application/octet-stream' });

          // Create a URL for the blob to be downloaded
          const url = URL.createObjectURL(blob);

          // Create a link element and download the file
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