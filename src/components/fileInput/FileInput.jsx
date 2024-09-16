import { useRef } from "react";
import { Flip, toast } from 'react-toastify';
import { FaFileArrowUp } from "react-icons/fa6"

const FileInput = ({ onFileChange }) => {
  const inputRef = useRef(null);

  const allowedFormats = ['.txt', '.csv']; // list of the allowed files to upload

    /* Toast invoke function */
    const notifyError = (message) => toast.error(message,{
      position: "bottom-right",
      transition: Flip,
      autoClose: 2500,
      closeOnClick: true,
      pauseOnHover: false,
      draggable: true
    });
  
  const handleDragOver = (e) => {
    e.preventDefault();
  }

  /* Handling file drop */
  const handleDrop = (e) => {
    e.preventDefault();
    onFileChange(e);

    const droppedFiles = e.dataTransfer.files;
    if (droppedFiles.length > 1) {
      notifyError("Only one file is allowed");
      return;
    }
    const file = droppedFiles[0];
    const fileExtension = file.name.split('.').pop().toLowerCase();
    if (!allowedFormats.includes(`.${fileExtension}`)) {
      notifyError(`File type .${fileExtension} is not allowed`);
      return;
    }
  }

  return (
    <div className="border-gray-600 border-2 border-dashed rounded-lg">
        <div className="flex flex-col items-center justify-center gap-2 p-8 relative">
          <FaFileArrowUp className="text-primary" size={38} />
          <p className="text-xl">Drag and drop a file here</p>
          <input
            accept=".txt"
            className="opacity-0 absolute inset-0 cursor-pointer"
            type="file"
            onDragOver={handleDragOver}
            onDrop={handleDrop} 
            ref={inputRef}
            onChange={onFileChange}/>
        </div> 
    </div>
  )
}

export default FileInput