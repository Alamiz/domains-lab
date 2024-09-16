import { useEffect, useRef, useState } from "react";
import { FileInput } from "../../components"
import { FaCloudArrowUp } from "react-icons/fa6";
import { FaFileLines } from "react-icons/fa6";
import ProgressBar from "../../components/progressBar/ProgressBar";
import { useFileUpload } from "../../hooks/useFileUpload";
import { Flip, ToastContainer, toast } from 'react-toastify';

const FileUpload = () => {
  const fileRef = useRef(null);
  const { file, uploading, processed, error, handleFileChange, progress } = useFileUpload();

  const notify = () => toast.success("File processed successfully !",{
    transition: Flip
  });

  useEffect(() => {
    if (processed)
      notify();
  }, [processed])

  return (
    <section>
      <div className="container">
        {!file ?
          <>
            <FileInput onFileChange={handleFileChange} />
            {/* Or pick a file */}
            <div className="flex items-center justify-center gap-4 mt-6">
              <p className="text-md font-bold">Or you can</p>
              <button className="flex items-center justify-center gap-2 text-background text-lg bg-primary rounded-full px-4 py-2"
                onClick={() => fileRef.current.click()} >
                Click here to upload <FaCloudArrowUp size={24} />
              </button>
            </div>
          </> :
          <div className="flex flex-col items-center justify-content">
            <FaFileLines size={56} className="text-primary mb-4" />
            <p className="text-xl mb-4">{file.name}</p>
            <ProgressBar progress={progress}/>

            {/* Upload another file */}
            <div className="flex items-center justify-center gap-4 mt-6">
              <p className="text-md font-bold">Upload another file ?</p>
              <button className="flex items-center justify-center gap-2 text-background text-lg bg-primary rounded-full px-4 py-2"
                onClick={() => fileRef.current.click()} >
                Click here to upload <FaCloudArrowUp size={24} />
              </button>
            </div>
          </div>
        }
        <input className="hidden" ref={fileRef} type="file" accept=".txt" onChange={handleFileChange} />
        <ToastContainer 
          position="bottom-right"
          autoClose={2500}
          closeOnClick
          pauseOnHover
          draggable />
      </div>
    </section>
  )
}

export default FileUpload