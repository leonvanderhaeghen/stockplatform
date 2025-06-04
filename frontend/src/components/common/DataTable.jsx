import React, { useMemo } from 'react';
import {
  DataGrid,
  GridToolbarContainer,
  GridToolbarExport,
  GridToolbarFilterButton,
  GridToolbarDensitySelector,
  GridToolbarColumnsButton,
  GridActionsCellItem,
} from '@mui/x-data-grid';
import { Box, Button, Typography, Chip, IconButton, Tooltip } from '@mui/material';
import { Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon, Visibility as ViewIcon } from '@mui/icons-material';

export const CustomToolbar = ({ title, onAdd, showAddButton = true, children }) => (
  <GridToolbarContainer sx={{ p: 2, justifyContent: 'space-between' }}>
    <Box>
      <Typography variant="h6" component="h2" sx={{ mb: 1 }}>
        {title}
      </Typography>
      <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
        <GridToolbarColumnsButton />
        <GridToolbarFilterButton />
        <GridToolbarDensitySelector />
        <GridToolbarExport
          printOptions={{
            hideFooter: true,
            hideToolbar: true,
          }}
        />
        {children}
      </Box>
    </Box>
    {showAddButton && onAdd && (
      <Button
        variant="contained"
        color="primary"
        startIcon={<AddIcon />}
        onClick={onAdd}
      >
        Add New
      </Button>
    )}
  </GridToolbarContainer>
);

const DataTable = ({
  rows = [],
  columns = [],
  loading = false,
  title = '',
  onAdd,
  onEdit,
  onDelete,
  onView,
  showAddButton = true,
  showActions = true,
  pageSize = 10,
  pageSizeOptions = [5, 10, 25],
  ...props
}) => {
  const tableColumns = useMemo(() => {
    const actionColumn = {
      field: 'actions',
      type: 'actions',
      headerName: 'Actions',
      width: 120,
      getActions: (params) => [
        onView && (
          <GridActionsCellItem
            key="view"
            icon={
              <Tooltip title="View">
                <ViewIcon />
              </Tooltip>
            }
            label="View"
            onClick={() => onView(params.row)}
            showInMenu={false}
          />
        ),
        onEdit && (
          <GridActionsCellItem
            key="edit"
            icon={
              <Tooltip title="Edit">
                <EditIcon color="primary" />
              </Tooltip>
            }
            label="Edit"
            onClick={() => onEdit(params.row)}
            showInMenu={false}
          />
        ),
        onDelete && (
          <GridActionsCellItem
            key="delete"
            icon={
              <Tooltip title="Delete">
                <DeleteIcon color="error" />
              </Tooltip>
            }
            label="Delete"
            onClick={() => onDelete(params.row)}
            showInMenu={false}
          />
        ),
      ].filter(Boolean),
    };

    return showActions ? [...columns, actionColumn] : columns;
  }, [columns, onEdit, onDelete, onView, showActions]);

  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <DataGrid
        rows={rows}
        columns={tableColumns}
        loading={loading}
        pageSizeOptions={pageSizeOptions}
        initialState={{
          pagination: {
            paginationModel: { page: 0, pageSize },
          },
        }}
        slots={{
          toolbar: () => (
            <CustomToolbar title={title} onAdd={onAdd} showAddButton={showAddButton}>
              {props.toolbarChildren}
            </CustomToolbar>
          ),
        }}
        slotProps={{
          toolbar: {
            showQuickFilter: true,
          },
        }}
        disableRowSelectionOnClick
        autoHeight
        {...props}
      />
    </Box>
  );
};

export default DataTable;
