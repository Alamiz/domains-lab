import { useEffect, useRef, useState } from "react";
import { FileInput } from "../../components"
import { FaCloudArrowUp } from "react-icons/fa6";
import { FaFileLines } from "react-icons/fa6";
import ProgressBar from "../../components/progressBar/ProgressBar";

const FileUpload = () => {
  const fileRef = useRef(null);
  const [file, setFile] = useState(null);

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
    e.target.value = null;
  }

  return (
    <section>
      <div className="container">
        {!file ?
          <>
            <FileInput file={file} setFile={setFile} />
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
            <ProgressBar progress={78}/>
          </div>
        }
        <input className="hidden" ref={fileRef} type="file" accept=".txt" onChange={handleFileChange} />
      </div>
    </section>
  )
}

export default FileUpload