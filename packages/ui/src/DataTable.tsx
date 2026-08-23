import type { ReactNode } from "react";

export type DataTableColumn<Row> = {
  id: string;
  header: ReactNode;
  cell: (row: Row) => ReactNode;
  headerClassName?: string;
  cellClassName?: string;
};

export type DataTableProps<Row> = {
  columns: Array<DataTableColumn<Row>>;
  rows: Array<Row>;
  rowKey: (row: Row) => string;
  emptyMessage?: string;
  caption?: string;
  className?: string;
  tableClassName?: string;
};

export function DataTable<Row>({
  columns,
  rows,
  rowKey,
  emptyMessage = "No hay datos para mostrar.",
  caption,
  className,
  tableClassName
}: DataTableProps<Row>) {
  return (
    <div className={`overflow-hidden rounded-t-[11px] border border-[#e1e4e8] bg-white ${className ?? ""}`.trim()}>
      <div className="overflow-x-auto">
        <table className={`w-full border-collapse ${tableClassName ?? ""}`.trim()}>
          {caption ? <caption className="sr-only">{caption}</caption> : null}
          <thead>
            <tr>
              {columns.map((column) => (
                <th
                  key={column.id}
                  className={`text-left ${column.headerClassName ?? ""}`.trim()}
                  scope="col"
                >
                  {column.header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {rows.length === 0 ? (
              <tr><td className="p-8 text-center text-[#64707b]" colSpan={columns.length}>{emptyMessage}</td></tr>
            ) : rows.map((row) => (
              <tr key={rowKey(row)}>
                {columns.map((column) => (
                  <td key={column.id} className={`text-left ${column.cellClassName ?? ""}`.trim()}>
                    {column.cell(row)}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
