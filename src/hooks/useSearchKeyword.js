import { useState } from 'react';
import axios from 'axios';

export const useSearchKeyword = () => {
    const [filepath, setFilepath] = useState(null);
    const [error, setError] = useState(null);

    const searchKeyword = async (keyword) => {
        try {
            const response = await axios.get(`${import.meta.env.VITE_DOMAINS_LAB_API}/search`, {
                params: { keyword },
            });
            setFilepath(response.data.filepath);
        } catch (error) {
            setError(error.message);
            console.log(error);
        }
    };

    return { filepath, error, searchKeyword };
};