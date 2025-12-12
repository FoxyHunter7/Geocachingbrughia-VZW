<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;

class ImageController
{
    public function get(string $p1 ,string $filename)
    {
        $path = storage_path("app/images/".$p1."/".$filename);

        if (!file_exists($path)) {
            return response()->json(["error" => "File not found"]);
        }

        return response()->file($path);
    }
}
