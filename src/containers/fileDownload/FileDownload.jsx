import { useState } from "react";
import { IoSearch } from "react-icons/io5";
import { FaDownload } from "react-icons/fa6";

const FileDownload = () => {
  const [searchKeyWord, setSearchKeyWord] = useState('')

  const search  = () => {
    console.log(searchKeyWord)
    setSearchKeyWord('')
  }

  return (
    <section>
      <div className="container">
        {/* Search bar */}
        <div className="flex gap-4">
          <input className="border-solid border border-gray-300 rounded px-4 py-2 w-full outline-gray-500" type="text" value={searchKeyWord} onChange={(e) => setSearchKeyWord(e.target.value)}/>
          <button className="bg-primary rounded px-3 py-3" onClick={search}>
            <IoSearch color="white" size={20}/>
          </button>
        </div>

        {/* Download results */}
        <div className="mt-8">
          <p className="text-3xl font-bold">Download results</p>
          <p className="mt-2">This is the result of your keyword Click to download.</p>
          <button className="flex bg-primary rounded-lg mt-3 px-3 py-3 text-white items-center justify-center gap-2">Download <FaDownload size={20}/></button>
        </div>
      </div>
    </section>
  )
}

export default FileDownload