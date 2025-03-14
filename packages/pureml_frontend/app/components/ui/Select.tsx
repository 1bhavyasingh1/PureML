import { tv, type VariantProps } from "tailwind-variants";
import * as SelectDropdown from "@radix-ui/react-select";
import { ChevronDown, MoreVertical } from "lucide-react";

const selectStyles = tv({
  slots: {
    base: "focus:outline-none",
    selectIcon:
      "focus:outline-none flex justify-center items-center text-base gap-x-2 font-medium text-slate-600 border border-slate-400 rounded px-4 py-2",
    noIcon: "focus:outline-none flex justify-center items-center",
    view: "bg-white flex justify-start rounded-lg shadow-lg",
  },
  variants: {
    intent: {
      primary: "text-slate-600",
      more: "text-slate-600",
    },
    fullWidth: {
      true: "w-full",
    },
  },
  defaultVariants: {
    intent: "primary",
    fullWidth: true,
  },
});
type SelectProps = VariantProps<typeof selectStyles>;

interface Props extends SelectProps {
  name: string;
  title: string;
  [children: string]: any;
}

export default function Select(props: Props) {
  const { base, selectIcon, noIcon, view } = selectStyles();
  return (
    <div>
      {props.intent === "primary" ? (
        <SelectDropdown.Root name={props.name}>
          <SelectDropdown.Trigger asChild className={base()}>
            <button className="flex foxus:outline-none rounded">
              <SelectDropdown.Icon className={selectIcon()}>
                {props.title} <ChevronDown className="w-4" />
              </SelectDropdown.Icon>
            </button>
          </SelectDropdown.Trigger>
          <SelectDropdown.Content
            position="popper"
            sideOffset={7}
            align="end"
            className="z-50"
          >
            <SelectDropdown.Viewport className={view()}>
              <SelectDropdown.Group className="min-w-[8rem]">
                {props.children}
              </SelectDropdown.Group>
            </SelectDropdown.Viewport>
          </SelectDropdown.Content>
        </SelectDropdown.Root>
      ) : (
        <SelectDropdown.Root name={props.name}>
          <SelectDropdown.Trigger asChild className={base()}>
            <button className="flex foxus:outline-none rounded">
              <SelectDropdown.Icon className={noIcon()}>
                <MoreVertical className="w-4" />
              </SelectDropdown.Icon>
            </button>
          </SelectDropdown.Trigger>
          <SelectDropdown.Content
            position="popper"
            sideOffset={7}
            align="end"
            className="z-50"
          >
            <SelectDropdown.Viewport className={view()}>
              <SelectDropdown.Group className="min-w-[8rem]">
                {props.children}
              </SelectDropdown.Group>
            </SelectDropdown.Viewport>
          </SelectDropdown.Content>
        </SelectDropdown.Root>
      )}
    </div>
  );
}
