package utils

// QC STATUS

var ASSET_QC_PENDING = 0
var ASSET_QC_COMPLETE = 1
var ASSET_PACKAGING_DONE = 2
var ASSET_IN_MOVEMENT = 3
var ASSET_RECEVIED_AT_WAREHOUSE = 4
var ASSET_DAMAGE_RECEIVED = 5
var ASSET_BUY_BACK = 6
var ASSET_SOLD = 7

// WORKFLOW QC STATUS
var QC_NOT_STARTED = 0
var QC_REQUESTED = 1
var QC_IN_PROGRESS = 2
var QC_COMPLETED = 3

// WORKFLOW SHIPMENT STATUS
var SHIPMENT_CREATED = 0
var PICKUP_REQUESTED = 1
var MOVEMENT_IN_PROGRESS = 2
var RECEIVED_AT_WAREHOUSE = 3

var SALES = "sales"
var ADMIN = "admin"
var TRANSFER = "transfer"
var QC = "qc"
var ENDOFTERM = "end_of_term"
var BSE = "bse"
var HR = "hr"

var PENDING = "pending"
var APPROVED = "approved"
var REJECTED = "rejected"
var PROCESSED = "processed"

// EPP
var ACTIVE = "active"
var INACTIVE = "inactive"
