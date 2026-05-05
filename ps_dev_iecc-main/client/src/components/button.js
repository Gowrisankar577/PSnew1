function CustomButton(props) {
  return (
    <button
      onClick={props.onClick}
      disabled={props.disabled}
      className={
        (props.danger ? "bg-red" : "bg-primary") +
        (props.smallFont ? " p-1 text-sm" : " p-2 text-md") +
        " w-full text-white tracking-wider rounded " +
        props.margin
      }
    >
      {props.label}
    </button>
  );
}

export default CustomButton;
