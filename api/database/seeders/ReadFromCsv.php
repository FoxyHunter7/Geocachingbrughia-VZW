<?php

namespace Database\Seeders;

class ReadFromCsv
{
    public static function getDataFromCsv($path, $separator = null)
    {
        $file = fopen($path,'r');
        $header = fgetcsv($file, separator:$separator);
        $data = [];
        while ($row = fgetcsv($file, separator:$separator)) {
            $data[] = array_combine($header, $row);
        }

        fclose($file);
        return $data;
    }
}
