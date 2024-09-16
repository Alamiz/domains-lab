import { useEffect, useState } from 'react';

export const useFileUpload = () => {
    const [file, setFile] = useState(null);
    const [uploading, setUploading] = useState(false);
    const [error, setError] = useState(null);
    const [progress, setProgress] = useState('');
    const [processed, setProcessed] = useState(false);

    /* If a file is selected, upload it */
    useEffect(() => {
        if (file)
            handleUpload();
    }, [file])

    const handleFileChange = (event) => {
        setFile(event.target.files[0]);
        event.target.value = '';
    };

    /* Upload file */
    const handleUpload = async () => {
        if (!file) return;
    
        setUploading(true);
        setProcessed(false);
        setProgress('0'); // Reset response data
    
        try {
            const formData = new FormData();
            formData.append('domainsFile', file); // creating a formData with the file to upload it
    
            const response = await fetch(`${import.meta.env.VITE_DOMAINS_LAB_API}/upload`, {
                method: 'POST',
                body: formData,
            });
    
            // Get a reader for the response body
            const reader = response.body.getReader();

            // Create a TextDecoder to decode the response body
            const decoder = new TextDecoder('utf-8');

            // Keep track of whether were done reading the response
            let done = false;

            // Loop until we've read the entire response
            while (!done) {
                // Read the next chunk of the response
                const { value, done: readerDone } = await reader.read();

                // Update our done flag
                done = readerDone;

                // Decode the chunk of the response we just read
                const chunk = decoder.decode(value, { stream: true });

                // Update the progress
                setProgress(chunk.split('\n')[0]);
            }

            // Set the progress to 100% and mark the upload as processed
            setProgress('100')
            setProcessed(true)
        } catch (error) {
            setError(error.message);
            console.log(error);
        } finally {
            setUploading(false);
        }
    };
    

    return { file, uploading, error, handleFileChange, handleUpload, progress, processed };
};