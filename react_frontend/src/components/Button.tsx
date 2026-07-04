type ButtonProps = {
    label: string
    click: () => void
}

function Button({ label, click }: ButtonProps) {
    return <button onClick={click}>{label}</button>
}

export default Button