import { useRef } from "react";
import { FaFileArrowUp } from "react-icons/fa6"
import { FaTrash } from "react-icons/fa6";

const FileInput = ({ onFileChange }) => {
  const inputRef = useRef(null);

  const handleDragOver = (e) => {
    e.preventDefault();
  }

  const handleDrop = (e) => {
    e.preventDefault();
    onFileChange(e);
  }


  return (
    <div className="border-gray-600 border-2 border-dashed rounded-lg">
      {/* {!file ?  */}
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
        {/* :  */}
        {/* <div className="flex items-center justify-center gap-2 p-8 relative">
          <p className="text-lg">{file.name}</p>
          <button>
            <FaTrash className="text-red-500"/>
          </button>
        </div> */}
        {/* } */}
    </div>
  )
}

export default FileInput