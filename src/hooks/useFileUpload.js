import { useEffect, useState } from 'react';
import axios from 'axios';

export const useFileUpload = () => {
    const [file, setFile] = useState(null);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState(null);
    const [progress, setProgress] = useState('');

    useEffect(() => {
        if (file)
            handleUpload();
    }, [file])

    const handleFileChange = (event) => {
        setFile(event.target.files[0]);
    };

    const handleUpload = async () => {
        if (!file) return;
    
        setUploading(true);
        setProgress(''); // Reset response data
    
        try {
            const formData = new FormData();
            formData.append('domainsFile', file);
    
            const response = await fetch(`${import.meta.env.VITE_DOMAINS_LAB_API}/upload`, {
                method: 'POST',
                body: formData,
            });
    
            const reader = response.body.getReader();
            const decoder = new TextDecoder('utf-8');
            let done = false;
    
            while (!done) {
                const { value, done: readerDone } = await reader.read();
                done = readerDone;
                const chunk = decoder.decode(value, { stream: true });
    
                // You can append the chunk to state to update it progressively
                setProgress(chunk.split('\n')[0]);
            }
    
            setProgress('100')
        } catch (error) {
            setError(error.message);
            console.log(error);
        } finally {
            setUploading(false);
        }
    };
    

    return { file, uploading, error, handleFileChange, handleUpload, progress };
};