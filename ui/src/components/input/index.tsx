const Input = ({
  type,
  className,
  placeholder,
  ...props
}: React.ComponentProps<"input">) => {
  return (
    <input
      type={type}
      className={className}
      placeholder={placeholder}
      {...props}
    />
  );
};

export default Input;
