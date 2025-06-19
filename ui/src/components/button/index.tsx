const Button = ({
  className,
  type,
  ...props
}: React.ComponentProps<"button">) => {
  return <button className={className} type={type} {...props} />;
};

export default Button;
