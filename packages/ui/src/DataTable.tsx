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
};

export function DataTable<Row>({
  columns,
  rows,
  rowKey,
  emptyMessage = "No hay datos para mostrar.",
  caption,
  className
}: DataTableProps<Row>) {
  return (
    <div className={`gs-data-table ${className ?? ""}`.trim()}>
      <div className="gs-data-table-scroll">
        <table>
          {caption ? <caption>{caption}</caption> : null}
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
              <tr><td className="gs-data-table-empty" colSpan={columns.length}>{emptyMessage}</td></tr>
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
