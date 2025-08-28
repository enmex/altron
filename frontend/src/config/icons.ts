import { IconType } from "react-icons";
import { BiFilterAlt, BiHeadphone, BiShareAlt, BiTime, BiTrash, BiBasket } from "react-icons/bi";
import { CgMergeVertical, CgOptions, CgSearch } from "react-icons/cg";
import { GiWarlockEye, GiDeathStar } from "react-icons/gi";
import { GoTerminal, GoTriangleDown, GoTriangleUp } from "react-icons/go";
import { IoMdAdd } from "react-icons/io";
import { MdArrowBackIosNew, MdArrowDropDown, MdContentCopy, MdDone, MdFileDownload, MdOutlineExpandMore, MdSignalWifiStatusbar2Bar, MdSignalWifiStatusbar4Bar, MdSignalWifiStatusbarNull, MdWarning } from "react-icons/md";
import { RiBrush3Line, RiEye2Line } from "react-icons/ri";
import { RxDash, RxReload } from "react-icons/rx";
import { SiConvertio } from "react-icons/si";
import { TbBrandPython, TbColorSwatch, TbHealthRecognition, TbScanEye, TbSquareLetterW } from "react-icons/tb";
import { VscDebugPause, VscDebugStart } from "react-icons/vsc";
import { HiArrowRight } from "react-icons/hi";
import { AiOutlineClose, AiOutlineSound } from "react-icons/ai";
import { CgServer } from "react-icons/cg";
import { FaMemory, FaDocker, FaUnlock, FaLock } from "react-icons/fa";
import { MdMemory, MdOutlineWorkspaces } from "react-icons/md";
import { CiLogout, CiImport } from "react-icons/ci";
import { BsFiletypeJson } from "react-icons/bs";

export interface Icons {
    burningEye: IconType;
    mdAdd: IconType;
    colorSwatch: IconType;
    arrowDropDown: IconType;
    scanEye: IconType;
    deathStar: IconType;
    options: IconType;
    debugPause: IconType;
    debugStart: IconType;
    eye2Line: IconType;
    brush3Line: IconType;
    terminal: IconType;
    letterW: IconType;
    outlineExpandMore: IconType;
    convertio: IconType;
    trash: IconType;
    headphones: IconType;
    outlineExport: IconType;
    search: IconType;
    filterAlt: IconType;
    arrowBackIosNew: IconType;
    contentCopy: IconType;
    reload: IconType;
    arrowRight: IconType;
    time: IconType;
    done: IconType;
    warning: IconType;
    download: IconType;
    share: IconType;
    sound: IconType;
    merge: IconType;
    server: IconType;
    health: IconType;
    memory: IconType;
    cpu: IconType;
    offline: IconType;
    online: IconType;
    mumble: IconType;
    increase: IconType; 
    decrease: IconType;
    stable: IconType;
    close: IconType;
    logout: IconType;
    basket: IconType;
    workspace: IconType;
    import: IconType;
    docker: IconType;
    lock: IconType;
    unlock: IconType;
    json: IconType;
}

export const icons: Icons = {
    burningEye: GiWarlockEye,
    mdAdd: IoMdAdd,
    colorSwatch: TbColorSwatch,
    arrowDropDown: MdArrowDropDown,
    scanEye: TbScanEye,
    deathStar: GiDeathStar,
    options: CgOptions,
    debugPause: VscDebugPause,
    debugStart: VscDebugStart,
    eye2Line: RiEye2Line,
    brush3Line: RiBrush3Line,
    terminal: GoTerminal,
    letterW: TbSquareLetterW,
    outlineExpandMore: MdOutlineExpandMore,
    convertio: SiConvertio,
    trash: BiTrash,
    headphones: BiHeadphone,
    outlineExport: TbBrandPython,
    search: CgSearch,
    filterAlt: BiFilterAlt,
    arrowBackIosNew: MdArrowBackIosNew,
    contentCopy: MdContentCopy,
    reload: RxReload,
    arrowRight: HiArrowRight,
    time: BiTime,
    done: MdDone,
    warning: MdWarning,
    download: MdFileDownload,
    share: BiShareAlt,
    sound: AiOutlineSound,
    merge: CgMergeVertical,
    server: CgServer,
    health: TbHealthRecognition,
    memory: FaMemory,
    cpu: MdMemory,
    offline: MdSignalWifiStatusbarNull,
    online: MdSignalWifiStatusbar4Bar,
    mumble: MdSignalWifiStatusbar2Bar,
    increase: GoTriangleUp,
    decrease: GoTriangleDown,
    stable: RxDash,
    close: AiOutlineClose,
    logout: CiLogout,
    basket: BiBasket,
    workspace: MdOutlineWorkspaces,
    import: CiImport,
    docker: FaDocker,
    lock: FaLock,
    unlock: FaUnlock,
    json: BsFiletypeJson
}