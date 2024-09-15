const ProgressBar = ({ progress }) => {
    return (
        <div className="flex flex-col justify-center w-full">
            <p className="mb-2">Processing file...</p>
            <div className="flex items-center gap-2 w-full">
                <div className="bg-violet-200 backdrop-blur h-3 rounded-full w-full">
                    <div className={`bg-primary h-full transition-all opacity-100 rounded-full`} style={{ width: `${progress}%` }}>
                    </div>
                </div>
                <p>{progress}%</p>
            </div>
        </div>
    )
}

export default ProgressBar