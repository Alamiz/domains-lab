const InstructionCard = ({title, description, icon}) => {
  return (
    <div className="p-6 rounded bg-primary brightness-110 border-2 border-white/20 text-white">
      {icon}
      <p className="text-lg font-bold">{title}</p>
      <p>{description}</p>
    </div>
  )
}

export default InstructionCard