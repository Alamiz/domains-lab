import { useRef, useState } from "react";
import { FileInput } from "../../components"
import { FaCloudArrowUp } from "react-icons/fa6";

const FileUpload = () => {
  const fileRef = useRef(null);
  const [file, setFile] = useState(null);

  return (
    <section>
      <div className="container">
        <FileInput file={file} setFile={setFile} />
        <div className="flex items-center justify-center gap-4 mt-6">
          <p className="text-md font-bold">Or you can</p>
          <button className="flex items-center justify-center gap-2 text-background text-lg bg-primary rounded-full px-4 py-2"
            onClick={() => fileRef.current.click()} >
            Click here to upload <FaCloudArrowUp size={24} />
          </button>
        </div>
        <input className="hidden" ref={fileRef} type="file" accept=".txt" onChange={(e) => setFile(e.target.files[0])} />
      </div>
    </section>
  )
}

export default FileUpload