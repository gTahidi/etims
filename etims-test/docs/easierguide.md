Okay, here's a more concise, developer-friendly, and professionally toned rewording of the provided document, maintaining all original meaning:

---

# eTIMS Online Sales Control Unit (OSCU) and Virtual Sales Control Unit (VSCU) Integration Guide

**Version 1.1**
**April 2023**
**ISO 9001:2015 CERTIFIED**

**Notice**

© 2023 Kenya Revenue Authority (KRA). This document is controlled. Unauthorized access, copying, replication, or usage for unintended purposes is prohibited. All trademarks are used for identification purposes only and belong to their respective companies.

---

This guide provides step-by-step instructions for integrating with the KRA's eTIMS system using either the Online Sales Control Unit (OSCU) or the Virtual Sales Control Unit (VSCU).

## A. eTIMS Taxpayer Sandbox Portal Sign-Up

1. **Access the Portal:** Navigate to the eTIMS taxpayer sandbox portal at [https://etims-sbx.kra.go.ke](https://etims-sbx.kra.go.ke).
2. **Initiate Sign-Up:** Click the "Sign Up" option.
3. **Complete the Form:**
    *   Select the "PIN" option.
    *   Enter your company's KRA PIN.
    *   Verify the pre-filled phone number. Ensure it's correct, as it will receive a One-Time Password (OTP). If incorrect, contact `timsupport@kra.go.ke` for an update.
    *   Create and confirm your password for the sandbox portal.

## B. OSCU or VSCU Device Registration

1. **Log In:** After successful signup, log in to the eTIMS taxpayer sandbox portal using your KRA PIN and password.
2. **Access Service Request:** Click the "Service Request" button located at the top right corner of the homepage.
3. **Select eTIMS:** In the dialog box, choose the "eTIMS" button.
4. **Complete Service Request Form:** Fill out the form, selecting the appropriate eTIMS Type:
    *   **VSCU (Virtual Sales Control Unit):** Client-side hosted solution suitable for high-volume, rapid invoice generation.
    *   **OSCU (Online Sales Control Unit):** KRA-hosted solution requiring integration with your invoicing system via the OSCU API. Ideal for online systems.
    *   **eTIMS Client:** KRA-developed standalone Windows and Android clients for PCs, mobile devices, PDAs, and tablets. Designed for Large and Medium Taxpayers (VAT and non-VAT registered), and other eligible applicants.
    *   **eTIMS Online:** Web-based invoicing via the eTIMS portal. Suitable for service sector businesses issuing fewer than 10 invoices monthly (e.g., rental property owners, transport providers, consultants, lawyers).
5. **eTIMS Commitment Form:** Complete and upload the required eTIMS Commitment Form, available on the Service Request form and the eTIMS "Learn More" page: [https://www.kra.go.ke/images/publications/eTIMS-confirmationdocument.pdf](https://www.kra.go.ke/images/publications/eTIMS-confirmationdocument.pdf).
6. **Confirmation:** Upon approval, you'll receive an SMS confirming your service request, allowing you to proceed with eTIMS installation.

## C. OSCU Setup

1. **Access and Activation:** Upon OSCU approval, access the OSCU hosted on KRA servers and begin the activation process.
2. **Refer to Documentation:** For detailed OSCU processes, configurations, and technical specifications, consult the OSCU Specifications document, section 2.2.
3. **API Endpoints:** Utilize the following base URLs for API interactions:
    *   **Sandbox:** `https://etims-api-sbx.kra.go.ke`
    *   **Production:** `https://etims-api.kra.go.ke`
    *   **Example:** For device activation (documented as `/selectInitOsdcInfo`), the full sandbox URL is `https://etims-api-sbx.kra.go.ke/selectInitOsdcInfo`.
4. **Initialization:** Your Trader Invoicing System (TIS) must invoke the OSCU initialization method, providing your PIN, branch office ID, and equipment information. The OSCU will then verify the device and retrieve a communication key from the KRA eTIMS API server.
5. **Crucial Credentials:** The Taxpayer PIN, branch office ID, and communication key are essential for subsequent TIS-OSCU communication.

## D. VSCU Setup

1. **Package Availability:** After VSCU approval, a deployment package (`eTIMS-VSCU-<version>.zip`) will be available on the eTIMS portal.
2. **Deployment:** Deploy the package to your server environment:
    *   **(a) Environment Setup:**
        *   Install Java Runtime Environment (JRE) or Java Development Kit (JDK) version 16 or higher.
        *   Configure Java environment variables (JAVA\_HOME and PATH).
    *   **(b) File Transfer:** Transfer the package to the target server using FTP, SCP, or a similar method. Then, unzip the package.
    *   **(c) Configuration:**
        *   Modify the `config/application.properties` file.
        *   The default port is 8088; change `server.port` if needed.
        *   Uncomment the appropriate `api.external.domain` for sandbox or production environments.
    *   **(d) Execution:**
        *   Open a command prompt or terminal.
        *   Navigate to the JAR file's directory.
        *   Run: `java -jar etims-vscu-<version>.jar`
    *   **(e) Connection:** Your TIS should connect to the VSCU server using: `http://<hostname or ip running vscu jar>:<server.port>`
3. **Refer to Documentation:** Consult the VSCU Specifications document, section 2.2, for detailed processes, configurations, and policies.
4. **API Endpoint:** Your TIS should interact with the VSCU server using the base URL: `http://<hostname or ip running vscu jar>:<server.port>`.
    *   **Example:** For device activation (documented as `/selectInitOsdcInfo`), if your VSCU server is at `vscuserverhostname`, the full URL is `http://vscuserverhostname:8088/selectInitOsdcInfo`.
5. **Initialization:** Your TIS must invoke the VSCU initialization service, providing your PIN, branch office ID, and equipment information in the request body. The VSCU will then verify the device, retrieve keys from the KRA eTIMS API server, and store them.

## E. Initialization and System Functionalities (Process Flow)

VSCU/OSCU functionalities are grouped into eight categories. These categories are sequential; later actions depend on earlier ones.

**1. Initialization (Send Only)**

*   **Purpose:** Maps your PIN, Branch Code, and serial number (provided during service request) to your TIS.
*   **Prerequisite:** eTIMS type registration and approval must be completed (refer to Section B).

**2. Basic Data Management (Get Only)**

*   **Purpose:** Retrieves essential standard codes and data from the eTIMS API Server for invoice generation. This includes item classification codes, location codes, package/weight codes, PIN lists, and KRA notices. Refer to the code definition chapter for details.
*   **Note:** Data consistency is crucial for sending invoice data.

**3. Branch Information Management (Get and Send)**

*   **Purpose:** Sends head/branch office lists and branch user information to the eTIMS API server.
*   **Pharmacy Note:** A function exists for sending insurance information.
*   **Note:** Branch codes are used for intra-branch stock transfers.

**4. Item Management (Get and Send)**

*   **Purpose:** Sends item information to and retrieves item lists from the eTIMS API Server.
*   **Note:** Allows item recovery from the eTIMS API server.

**5. Imported Item Management (Get and Send)**

*   **Purpose:** Retrieves imported item data (declared under your TIS owner's PIN) from the KRA customs system. Sends confirmation of received items with corresponding TIS stock items.
*   **Note:** Imported item data can be used for stock adjustments.

**6. Sales Management (Send Only)**

*   **Purpose:** Sends sales transaction and invoice data to the eTIMS API Server.
*   **VSCU Note:** For VSCU sales transaction information must proceed sales Invoice information.
*   **Sales Transaction Data:** Includes Customer PIN, Customer Name, Sales Type Code, Receipt Type Code, Payment Type Code, Invoice Status Code, Validated Date, Sale Date, Stock Released Date, Cancel Requested Date, Canceled Date, and Refunded Date.
*   **Sales Invoice Data:** Includes Invoice Number, Current Receipt Number, Total Receipt Number, Customer PIN, Customer Mobile Number, Receipt Published Date, Internal Data, and Receipt Signature.
* **VSCU Only Note**: TIS application sends sales invoice information to the VSCU and gets Receipt counters, Receipt Date, Internal Data, and Signature Data.

**7. Purchase Transaction Management (Get and Send)**

*   **Purpose:** Retrieves purchase transactions and invoice data from the eTIMS Server (under your TIS owner's PIN). Allows confirmation of purchases for stock adjustment.

**8. Stock Management (Get and Send)**

*   **Purpose:** Sends inventory in/out data for branches and updates stock status by item classification. Allows requesting stock from the main branch.
*   **Note**: Every stock in/out information must have its sales Invoice information sent in advance to eTIMS Server.
* **Note**: Every stock inventory information must have its stock in/out information sent in advance to eTIMS Server.

---

**Functionality Summary Table**

| #   | Category                      | Action of TIS side                                             | Description                                                                                                                                                                                        | Remark                                    |
| :-- | :---------------------------- | :------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :---------------------------------------- |
| 1   | Initialization                | Device authentication                                          | Device authentication from KRA. 3rd party users must complete an eTIMS service request.                                                                                                          |                                           |
| 2   | Basic data management         | Get code list                                                 | TIS application can update common standard codes managed by KRA from eTIMS API server                                                                                                           |                                           |
|     |                               | Get item classification list                                  | Server TIS application can update item classification codes managed by KRA from eTIMS API server.                                                                                                 |                                           |
|     |                               | Get PIN information                                           | TIS application can get information on a certain PIN from eTIMS server                                                                                                                            |                                           |
|     |                               | Get the branch list of head office(store)                    | TIS application can update the branch office information of head office into eTIMS API server.                                                                                                  |                                           |
|     |                               | Get notice list from eTIMS API server                        | TIS application gets eTIMS notification managed by KRA from eTIMS API server.                                                                                                                  |                                           |
| 3   | Branch information management | Send customer information                                     | TIS application sends customer information of the head & branch offices into eTIMS Server.                                                                                                      |                                           |
|     |                               | Send branch user account                                      | TIS application sends branch user account information to upload it in eTIMS Server.                                                                                                             |                                           |
|     |                               | Send branch insurance information                            | TIS application sends insurance information of the head & branch offices and update it in eTIMS Server                                                                                            | Applied to Pharmacy only. Not Mandatory |
| 4   | Item management               | Send Item information                                        | TIS application sends item information (name, unit price, bar code and etc.) of the head or branch offices and updates them in eTIMS Server                                                         | Head office(store) role                 |
|     |                               | Send Item Composition information                            | TIS application sends item composition information of the head or branch store and updates them in eTIMS API server. Item composition is used to calculate BOM (Bill of Material) of the final processed items. | Head office(store) role                 |
|     |                               | Get item list                                                | TIS application gets item information updated to eTIMS Server by TIS application.                                                                                                              |                                           |
| 5   | Imported item management      | Get imported item information                                | TIS application gets imported item information managed by KRA customs department. All the import declaration result of an owner of TIS application is provided from eTIMS Server.                 | Head office(store) role                 |
|     |                               | Send (converted) imported item information                   | TIS shall convert the received imported item information into eTIMS standard item information for sales (inventory). TIS application sends converted imported items information to eTIMS API server. |                                           |
| 6   | Sales management              | Send sales transaction and sales Invoice information         | TIS application sends sales transaction and sales invoice information to eTIMS API server. The sales transaction information must proceed sales Invoice information for VSCU.             |                                           |
|     |                               | Send sales Invoice information and get Internal Data and Signature Data | TIS application sends sales invoice information to the VSCU and gets Receipt counters, Receipt Date, Internal Data and Signature Data. Every sales Invoice information must have its sales transaction information sent in advance to eTIMS server.          | VSCU Only |
| 7   | Purchase transaction management | Get purchase transaction information                        | TIS application gets purchase transaction information which were registered to eTIMS Server by a seller under TIS application owner’s PIN.                                                           | Head office(store) role                 |
|     |                               | Send purchase transaction confirmation                       | TIS application sends purchase transaction confirmation of received purchase transaction information to eTIMS Server.                                                                             | Head office(store) role                 |
| 8   | Stock management              | Send stock in/out information                                | When adjusting inventory quantities or selling items at each branches or head office, TIS application sends the stock in/out information to the eTIMS Server.                                      |                                           |
|     |                               | Send stock inventory information                             | If stock inventory information is changed due to stock IN / OUT, Server TIS application sends the changed information of the inventory master and stores it in the eTIMS API server.            |                                           |
|     |                               | Receive stock from head office.                              | Branch office should request for inventory coming from head office.                                                                                                                              |                                           |

---
Let me know if you have any further refinements!
