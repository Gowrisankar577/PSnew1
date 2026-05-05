function InputBox(props) {
  return (
    <div className={"w-full " + (props.margin || "")}>
      <label style={{ fontSize: 14 }}>{props.label}</label>
      <input
        onChange={(e) => {
          props.return === "target"
            ? props.onChange(e)
            : props.type === "file"
              ? props.onChange(e.target.files[0])
              : props.onChange(e.target.value);
        }}
        style={{
          padding: 6,
          paddingLeft: 10,
          paddingRight: 10,
          backgroundColor: props.backgroundColour || "rgb(238 241 249/1)",
        }}
        value={props.value}
        disabled={props.disabled}
        className="border w-full mt-1 rounded placeholder:text-sm focus:outline-primary disabled:opacity-50 disabled:cursor-not-allowed"
        placeholder={props.placeholder}
        type={props.type}
        accept={props.accept}
      />
    </div>
  );
}

export default InputBox;
