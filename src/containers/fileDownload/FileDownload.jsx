import { useEffect, useState } from "react";
import { IoSearch } from "react-icons/io5";
import { RiLoader3Fill } from "react-icons/ri";
import { FaDownload } from "react-icons/fa6";
import { useSearchKeyword } from "../../hooks/useSearchKeyword";

const FileDownload = ({ isFileProcessed }) => {
  const { filepath, error, searchKeyword, downloadFile, loading } = useSearchKeyword()

  const [keyWord, setKeyWord] = useState('')
  const [inputError, setInputError] = useState('')

  const search  = () => {
    if (!keyWord) {
      setInputError('Please enter a keyword')
      return
    }
    searchKeyword(keyWord)
    setKeyWord('')
  }

  useEffect(() => {
    if (keyWord) {
      setInputError('')
    }
  }, [keyWord])

  return (
    <section>
      <div className="container">
        {/* Search bar */}
        <div className="flex flex-col">
          {inputError && <p className="text-sm text-red-500">{inputError}</p>}
          <div className="flex gap-4 w-full">
            <input className={`${inputError ? 'border-red-500' : ''} ${!isFileProcessed || loading ? 'bg-gray-200' : ''} border-solid border border-gray-300 rounded px-4 py-2 w-full outline-gray-500`} disabled={!isFileProcessed || loading} type="text" value={keyWord} onChange={(e) => setKeyWord(e.target.value)}/>
            <button className="bg-primary rounded px-3 py-3" onClick={search} disabled={!isFileProcessed || loading}>
              {
                loading ?
                  <RiLoader3Fill className="animate-spin" color="white" size={24} /> :
                  <IoSearch color="white" size={24} />
              }
            </button>
          </div>
        </div>

        {/* Download results */}
        {(filepath && !loading) && <div className="mt-8">
          <p className="text-3xl font-bold">Download results</p>
          <p className="mt-2">This is the result of your keyword Click to download.</p>
          <button className="flex bg-primary rounded-lg mt-3 px-3 py-3 text-white items-center justify-center gap-2" onClick={() => downloadFile(filepath)}>Download <FaDownload size={20} /></button>
        </div>}
      </div>
    </section>
  )
}

export default FileDownload