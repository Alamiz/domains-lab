import { useEffect, useRef } from "react";
import { useFileUpload } from "../../hooks/useFileUpload";
import ProgressBar from "../../components/progressBar/ProgressBar";
import { Flip, toast } from 'react-toastify';
import { FileInput } from "../../components"
import { FaCloudArrowUp } from "react-icons/fa6";
import { FaFileLines } from "react-icons/fa6";

const FileUpload = () => {
  const fileRef = useRef(null);
  const { file, processed, error, handleFileChange, progress } = useFileUpload();

  /* Toast invoke function */
  const notify = () => toast.success("File processed successfully !", {
    position: "bottom-right",
    transition: Flip,
    autoClose: 2500,
    closeOnClick: true,
    pauseOnHover: false,
    draggable: true
  });

  /* Toast invoke function */
  const notifyError = () => toast.error(error, {
    position: "bottom-right",
    transition: Flip,
    autoClose: 2500,
    closeOnClick: true,
    pauseOnHover: false,
    draggable: true
  });

  /* Notify on success */
  useEffect(() => {
    if (processed) notify();
    if (error) notifyError();
  }, [processed, error])

  return (
    <section id="upload">
      <div className="container">
        {!file ?
          <>
            <FileInput onFileChange={handleFileChange} />
            {/* Or pick a file */}
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mt-6">
              <p className="text-md font-bold">Or you can</p>
              <button className="flex items-center justify-center gap-2 text-background text-lg bg-primary rounded-full px-4 py-2"
                onClick={() => fileRef.current.click()} >
                Click here to upload <FaCloudArrowUp size={24} />
              </button>
            </div>
          </> :
          <div className="flex flex-col items-center justify-content">
            {/* File preview */}
            <FaFileLines size={56} className="text-primary mb-4" />
            <p className="text-xl mb-4">{file.name}</p>
            <ProgressBar progress={progress} />

            {/* Upload another file */}
            {processed && <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mt-6">
              <p className="text-md font-bold">Upload another file ?</p>
              <button className="flex items-center justify-center gap-2 text-background text-lg bg-primary rounded-full px-4 py-2"
                onClick={() => fileRef.current.click()} >
                Click here to upload <FaCloudArrowUp size={24} />
              </button>
            </div>}
          </div>
        }
        <input className="hidden" ref={fileRef} type="file" accept=".txt" onChange={handleFileChange} />
      </div>
    </section>
  )
}

export default FileUpload