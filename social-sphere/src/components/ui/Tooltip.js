export default function Tooltip({ content, active=true, children }) {

    if (active === false) {
        return (
            <div className="group/tooltip relative inline-flex">
            {children}
        </div>
        )
    }

    return (
        <div className="group/tooltip relative inline-flex">
            {children}
            <div className="absolute top-full left-1/2 -translate-x-1/2 mt-2 px-2.5 py-1 bg-(--accent) text-white text-xs font-medium rounded-md opacity-0 group-hover/tooltip:opacity-100 transition-opacity duration-200 whitespace-nowrap pointer-events-none z-50 shadow-lg">
                {content}
                {/* Tiny arrow pointing up */}
                <div className="absolute -top-1 left-1/2 -translate-x-1/2 w-2 h-2 bg-(--accent) rotate-45" />
            </div>
        </div>
    );
}